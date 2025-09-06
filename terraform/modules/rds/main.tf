# terraform/modules/rds/main.tf

# DBがどのサブネットに配置されるかを定義するグループ
resource "aws_db_subnet_group" "main" {
  name       = "${var.project_name}-subnet-group"
  subnet_ids = var.private_subnet_ids
}

# DBへのアクセスを制御するセキュリティグループ（ファイアウォール）
resource "aws_security_group" "db" {
  name   = "${var.project_name}-db-sg"
  vpc_id = var.vpc_id
  # ingress { # TODO: 後でECSからのアクセスを許可するルールを追加 }
}

# PostgreSQLデータベース本体
resource "aws_db_instance" "main" {
  identifier           = "${var.project_name}-db"
  engine               = "postgres"
  instance_class       = "db.t3.micro" # MVP用の最小インスタンス
  allocated_storage    = 20
  username             = var.db_username
  password             = var.db_password
  db_subnet_group_name = aws_db_subnet_group.main.name
  vpc_security_group_ids = [aws_security_group.db.id]
  skip_final_snapshot  = true
}