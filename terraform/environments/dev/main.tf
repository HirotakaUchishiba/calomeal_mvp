# terraform/environments/dev/main.tf

terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

# AWS Provider設定
provider "aws" {
  region = var.aws_region
  
  default_tags {
    tags = {
      Environment = "dev"
      Project     = "calomeal"
      ManagedBy   = "terraform"
    }
  }
}

# ローカル変数
locals {
  project_name = "calomeal"
  environment  = "dev"
}

# VPCモジュール
module "vpc" {
  source = "../../modules/vpc"
  
  project_name = local.project_name
  environment  = local.environment
  cidr_block   = "10.0.0.0/16"
  
  availability_zones = ["${var.aws_region}a", "${var.aws_region}b"]
  
  public_subnet_cidrs  = ["10.0.1.0/24", "10.0.2.0/24"]
  private_subnet_cidrs = ["10.0.10.0/24", "10.0.20.0/24"]
  
  tags = {
    Environment = local.environment
    Project     = local.project_name
  }
}

# RDSモジュール
module "rds" {
  source = "../../modules/rds"
  
  project_name = local.project_name
  environment  = local.environment
  
  vpc_id             = module.vpc.vpc_id
  private_subnet_ids = module.vpc.private_subnet_ids
  security_group_ids = [module.vpc.default_security_group_id]
  
  db_instance_class = "db.t3.micro"
  db_allocated_storage = 20
  db_engine_version = "15.4"
  
  tags = {
    Environment = local.environment
    Project     = local.project_name
  }
}

# Cognitoモジュール
module "cognito" {
  source = "../../modules/cognito"
  
  project_name = local.project_name
  environment  = local.environment
  
  # 開発環境用のコールバックURL
  callback_urls = [
    "http://localhost:5173",
    "https://localhost:5173",
    "http://localhost:3000",
    "https://localhost:3000"
  ]
  
  logout_urls = [
    "http://localhost:5173",
    "https://localhost:5173",
    "http://localhost:3000",
    "https://localhost:3000"
  ]
  
  # 開発環境ではIdentity Poolは無効
  enable_identity_pool = false
  
  tags = {
    Environment = local.environment
    Project     = local.project_name
  }
}

# IAMモジュール
module "iam" {
  source = "../../modules/iam"
  
  project_name = local.project_name
  environment  = local.environment
  
  tags = {
    Environment = local.environment
    Project     = local.project_name
  }
}

# ECSモジュール
module "ecs" {
  source = "../../modules/ecs"
  
  project_name = local.project_name
  environment  = local.environment
  
  vpc_id             = module.vpc.vpc_id
  public_subnet_ids  = module.vpc.public_subnet_ids
  private_subnet_ids = module.vpc.private_subnet_ids
  
  # データベース接続情報
  db_endpoint = module.rds.db_endpoint
  db_name     = module.rds.db_name
  db_username = module.rds.db_username
  db_password = module.rds.db_password
  
  # Cognito設定
  cognito_user_pool_id     = module.cognito.user_pool_id
  cognito_client_id        = module.cognito.user_pool_client_id
  cognito_user_pool_domain = module.cognito.user_pool_domain
  
  tags = {
    Environment = local.environment
    Project     = local.project_name
  }
}
