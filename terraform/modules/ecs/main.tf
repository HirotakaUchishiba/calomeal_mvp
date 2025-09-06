# terraform/modules/ecs/main.tf

# 現在のAWSアカウント情報を取得
data "aws_caller_identity" "current" {}

# -----------------------------------------------------------------------------
# 変数定義 (variables.tf に記述することを推奨)
# このモジュールが外部から受け取る設定値
# -----------------------------------------------------------------------------

variable "project_name" {
  description = "プロジェクト名。リソースの命名に使用"
  type        = string
}

variable "environment" {
  description = "環境名 (例: dev, prd)"
  type        = string
}

variable "vpc_id" {
  description = "ECSサービスを配置するVPCのID"
  type        = string
}

variable "private_subnet_ids" {
  description = "ECSタスクを配置するプライベートサブネットのIDリスト"
  type        = list(string)
}

variable "alb_target_group_arn" {
  description = "トラフィックを受け取るALBターゲットグループのARN"
  type        = string
}

variable "container_image" {
  description = "デプロイするDockerコンテナイメージのURI (ECRリポジトリURI)"
  type        = string
}

variable "container_port" {
  description = "コンテナがリッスンするポート番号"
  type        = number
  default     = 8080
}

variable "db_secret_arn" {
  description = "データベースの認証情報を格納したAWS Secrets ManagerのシークレットARN"
  type        = string
}

# -----------------------------------------------------------------------------
# ECSクラスタ
# ECSサービスとタスクの論理的なグループ。
# -----------------------------------------------------------------------------
resource "aws_ecs_cluster" "this" {
  name = "${var.project_name}-${var.environment}-cluster"

  tags = {
    Name        = "${var.project_name}-${var.environment}-cluster"
    Project     = var.project_name
    Environment = var.environment
  }
}

# -----------------------------------------------------------------------------
# IAMロールとポリシー
# 設計資料の指示通り、タスク実行ロールとタスクロールを分離し、最小権限の原則に従う。
# -----------------------------------------------------------------------------

# 1. ECSタスク実行ロール (Task Execution Role)
# ECSエージェントがECRからDockerイメージをプルしたり、CloudWatchにログを送信したりするために必要な権限。
resource "aws_iam_role" "ecs_task_execution_role" {
  name = "${var.project_name}-${var.environment}-ecs-task-execution-role"

  assume_role_policy = jsonencode({
    Version   = "2012-10-17"
    Statement = [
  {
    Effect = "Allow"
    Principal = {
      Service = "ecs-tasks.amazonaws.com"
    }
    Action = "sts:AssumeRole"
  }
]
  })

  tags = {
    Project     = var.project_name
    Environment = var.environment
  }
}

# AWS管理ポリシーをアタッチ
resource "aws_iam_role_policy_attachment" "ecs_task_execution_role_policy" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# 2. ECSタスクロール (Task Role)
# コンテナ内のアプリケーション自体が、他のAWSサービス(Secrets Managerなど)にアクセスするために必要な権限。
resource "aws_iam_role" "ecs_task_role" {
  name = "${var.project_name}-${var.environment}-ecs-task-role"

  assume_role_policy = jsonencode({
    Version   = "2012-10-17"
     Statement = [
  {
    Effect = "Allow"
    Principal = {
      Service = "ecs-tasks.amazonaws.com"
    }
    Action = "sts:AssumeRole"
  }
]
  })

  tags = {
    Project     = var.project_name
    Environment = var.environment
  }
}

# アプリケーションがSecrets ManagerからDB認証情報を読み取るためのカスタムポリシー
resource "aws_iam_policy" "secrets_manager_read_policy" {
  name        = "${var.project_name}-${var.environment}-secrets-manager-read-policy"
  description = "Allows reading specific secrets from AWS Secrets Manager"

  policy = jsonencode({
    Version   = "2012-10-17"
    Statement = [
  {
    Effect = "Allow"
    Action = [
      "secretsmanager:GetSecretValue"
    ]
    Resource = "arn:aws:secretsmanager:*:*:secret:${var.project_name}-${var.environment}-db-credentials-*"
  }
]
  })
}

resource "aws_iam_role_policy_attachment" "ecs_task_role_secrets_manager_policy" {
  role       = aws_iam_role.ecs_task_role.name
  policy_arn = aws_iam_policy.secrets_manager_read_policy.arn
}


# -----------------------------------------------------------------------------
# ECSタスク定義
# 実行するコンテナの設計図。CPU/メモリ、イメージ、ポートマッピング、環境変数、シークレットなどを定義。
# -----------------------------------------------------------------------------
resource "aws_ecs_task_definition" "this" {
  family                   = "${var.project_name}-${var.environment}-task"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"  # 0.25 vCPU
  memory                   = "512"  # 512 MiB
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn

  # コンテナ定義
  container_definitions = jsonencode([
    {
      name      = "${var.project_name}-container"
      image     = var.container_image
      essential = true
      portMappings = [
        {
          containerPort = var.container_port
          hostPort      = var.container_port
          protocol      = "tcp"
        }
      ]
      # 設計資料の指示通り、AWS Secrets ManagerからDB認証情報を安全に注入する。
      # これにより、機密情報がコードやコンテナイメージにハードコードされるのを防ぐ。
      secrets = [
  {
    name      = "DATABASE_URL"
    valueFrom = "arn:aws:secretsmanager:ap-northeast-1:${data.aws_caller_identity.current.account_id}:secret:${var.project_name}-${var.environment}-db-credentials"
  }
]
      # 構造化JSON形式でログをCloudWatch Logsに出力
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = "/ecs/${var.project_name}-${var.environment}"
          "awslogs-region"        = "ap-northeast-1" # リージョンは適宜変更
          "awslogs-stream-prefix" = "ecs"
        }
      }
    }
  ])

  tags = {
    Project     = var.project_name
    Environment = var.environment
  }
}

# -----------------------------------------------------------------------------
# CloudWatch Logs グループ
# コンテナからのログを収集・保管する場所。
# -----------------------------------------------------------------------------
resource "aws_cloudwatch_log_group" "this" {
  name              = "/ecs/${var.project_name}-${var.environment}"
  retention_in_days = 7 # ログの保持期間

  tags = {
    Project     = var.project_name
    Environment = var.environment
  }
}

# -----------------------------------------------------------------------------
# ECSサービス
# タスク定義に基づき、指定された数のタスクを常に実行し続ける責務を持つ。
# ALBとの連携やネットワーク設定もここで行う。
# -----------------------------------------------------------------------------
resource "aws_ecs_service" "this" {
  name            = "${var.project_name}-${var.environment}-service"
  cluster         = aws_ecs_cluster.this.id
  task_definition = aws_ecs_task_definition.this.arn
  desired_count   = 1 # MVPでは1つのタスクを実行
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = var.private_subnet_ids
    # セキュリティグループはVPCモジュール等で作成し、変数で渡すのが一般的
    # security_groups = [var.ecs_security_group_id]
    assign_public_ip = false # プライベートサブネットに配置するため、パブリックIPは不要
  }

  load_balancer {
    target_group_arn = var.alb_target_group_arn
    container_name   = "${var.project_name}-container"
    container_port   = var.container_port
  }

  # 新しいデプロイメントが完了する前に古いタスクを停止しないようにする設定
  # これにより、デプロイ中のダウンタイムを防ぐ
  deployment_minimum_healthy_percent = 100
  deployment_maximum_percent         = 200

  # ALBからのヘルスチェックが失敗した場合に、ECSがタスクを自動的に置き換えるのを待つ時間
  health_check_grace_period_seconds = 60

  # aws_ecs_cluster.thisへの明示的な依存関係
  depends_on = [aws_iam_role_policy_attachment.ecs_task_execution_role_policy]

  tags = {
    Project     = var.project_name
    Environment = var.environment
  }
}