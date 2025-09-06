# terraform/environments/prd/main.tf

provider "aws" {
  region = "ap-northeast-1" # 東京リージョン
}

# 現在のAWSアカウント情報を取得
data "aws_caller_identity" "current" {}

# --- VPCモジュールの呼び出し ---
module "vpc" {
  source = "../../modules/vpc"

  project_name         = "calomeal-prd"
  region               = "ap-northeast-1"
  vpc_cidr             = "10.0.0.0/16"
  public_subnets_cidr  = ["10.0.1.0/24", "10.0.2.0/24"]
  private_subnets_cidr = ["10.0.101.0/24", "10.0.102.0/24"]
}

# --- IAMモジュールの呼び出し ---
module "iam" {
  source = "../../modules/iam"
}

# --- RDSモジュールの呼び出し ---
module "rds" {
  source = "../../modules/rds"

  project_name       = "calomeal-prd"
  vpc_id             = module.vpc.vpc_id
  private_subnet_ids = module.vpc.private_subnet_ids
  db_username        = "admin"
  db_password        = var.db_password # 変数で安全に渡す
}

# --- ECSモジュールの呼び出し ---
module "ecs" {
  source = "../../modules/ecs"

  project_name              = "calomeal-prd"
  environment               = "prd"
  vpc_id                    = module.vpc.vpc_id
  private_subnet_ids        = module.vpc.private_subnet_ids
  container_image           = "your-aws-account-id.dkr.ecr.ap-northeast-1.amazonaws.com/calomeal:latest" # 仮
  db_secret_arn            = "arn:aws:secretsmanager:ap-northeast-1:${data.aws_caller_identity.current.account_id}:secret:calomeal-prd-db-credentials"
  alb_target_group_arn     = "arn:aws:elasticloadbalancing:ap-northeast-1:${data.aws_caller_identity.current.account_id}:targetgroup/calomeal-prd-tg"
}