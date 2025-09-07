# terraform/environments/prd/main.tf

terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "ap-northeast-1" # 東京リージョン
  
  default_tags {
    tags = {
      Environment = "prd"
      Project     = "calomeal"
      ManagedBy   = "terraform"
    }
  }
}

# 現在のAWSアカウント情報を取得
data "aws_caller_identity" "current" {}

# ローカル変数
locals {
  project_name = "calomeal"
  environment  = "prd"
}

# --- VPCモジュールの呼び出し ---
module "vpc" {
  source = "../../modules/vpc"

  project_name         = local.project_name
  environment          = local.environment
  region               = "ap-northeast-1"
  vpc_cidr             = "10.0.0.0/16"
  public_subnets_cidr  = ["10.0.1.0/24", "10.0.2.0/24"]
  private_subnets_cidr = ["10.0.101.0/24", "10.0.102.0/24"]
  
  tags = {
    Environment = local.environment
    Project     = local.project_name
  }
}

# --- IAMモジュールの呼び出し ---
module "iam" {
  source = "../../modules/iam"
  
  project_name = local.project_name
  environment  = local.environment
  
  tags = {
    Environment = local.environment
    Project     = local.project_name
  }
}

# --- RDSモジュールの呼び出し ---
module "rds" {
  source = "../../modules/rds"

  project_name       = local.project_name
  environment        = local.environment
  vpc_id             = module.vpc.vpc_id
  private_subnet_ids = module.vpc.private_subnet_ids
  db_username        = "admin"
  db_password        = var.db_password # 変数で安全に渡す
  
  # 本番環境用の設定
  db_instance_class    = "db.t3.small"
  db_allocated_storage = 100
  db_engine_version    = "15.4"
  
  tags = {
    Environment = local.environment
    Project     = local.project_name
  }
}

# --- Cognitoモジュールの呼び出し ---
module "cognito" {
  source = "../../modules/cognito"
  
  project_name = local.project_name
  environment  = local.environment
  
  # 本番環境用のコールバックURL（実際のドメインに変更）
  callback_urls = [
    "https://calomeal.com",
    "https://www.calomeal.com",
    "https://app.calomeal.com"
  ]
  
  logout_urls = [
    "https://calomeal.com",
    "https://www.calomeal.com",
    "https://app.calomeal.com"
  ]
  
  # 本番環境ではIdentity Poolを有効
  enable_identity_pool = true
  
  # 本番環境用のカスタムドメイン（オプション）
  custom_domain = var.cognito_custom_domain
  
  tags = {
    Environment = local.environment
    Project     = local.project_name
  }
}

# --- ECSモジュールの呼び出し ---
module "ecs" {
  source = "../../modules/ecs"

  project_name              = local.project_name
  environment               = local.environment
  vpc_id                    = module.vpc.vpc_id
  private_subnet_ids        = module.vpc.private_subnet_ids
  container_image           = "your-aws-account-id.dkr.ecr.ap-northeast-1.amazonaws.com/calomeal:latest" # 仮
  db_secret_arn            = "arn:aws:secretsmanager:ap-northeast-1:${data.aws_caller_identity.current.account_id}:secret:calomeal-prd-db-credentials"
  alb_target_group_arn     = "arn:aws:elasticloadbalancing:ap-northeast-1:${data.aws_caller_identity.current.account_id}:targetgroup/calomeal-prd-tg"
  
  # Cognito設定
  cognito_user_pool_id     = module.cognito.user_pool_id
  cognito_client_id        = module.cognito.user_pool_client_id
  cognito_user_pool_domain = module.cognito.user_pool_domain
  
  tags = {
    Environment = local.environment
    Project     = local.project_name
  }
}