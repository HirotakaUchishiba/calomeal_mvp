# terraform/environments/dev/outputs.tf

# VPC出力
output "vpc_id" {
  description = "ID of the VPC"
  value       = module.vpc.vpc_id
}

output "public_subnet_ids" {
  description = "IDs of the public subnets"
  value       = module.vpc.public_subnet_ids
}

output "private_subnet_ids" {
  description = "IDs of the private subnets"
  value       = module.vpc.private_subnet_ids
}

# RDS出力
output "db_endpoint" {
  description = "RDS instance endpoint"
  value       = module.rds.db_endpoint
  sensitive   = true
}

output "db_name" {
  description = "RDS database name"
  value       = module.rds.db_name
}

output "db_username" {
  description = "RDS master username"
  value       = module.rds.db_username
}

output "db_password" {
  description = "RDS master password"
  value       = module.rds.db_password
  sensitive   = true
}

# Cognito出力
output "cognito_user_pool_id" {
  description = "Cognito User Pool ID"
  value       = module.cognito.user_pool_id
}

output "cognito_client_id" {
  description = "Cognito User Pool Client ID"
  value       = module.cognito.user_pool_client_id
}

output "cognito_user_pool_domain" {
  description = "Cognito User Pool Domain"
  value       = module.cognito.user_pool_domain
}

output "cognito_config" {
  description = "Cognito configuration for environment variables"
  value       = module.cognito.cognito_config
}

# ECS出力
output "ecs_cluster_name" {
  description = "Name of the ECS cluster"
  value       = module.ecs.cluster_name
}

output "ecs_service_name" {
  description = "Name of the ECS service"
  value       = module.ecs.service_name
}

output "alb_dns_name" {
  description = "DNS name of the Application Load Balancer"
  value       = module.ecs.alb_dns_name
}

# 環境変数用の統合出力
output "environment_variables" {
  description = "Environment variables for the application"
  value = {
    # データベース設定
    DB_HOST     = module.rds.db_endpoint
    DB_PORT     = "5432"
    DB_NAME     = module.rds.db_name
    DB_USERNAME = module.rds.db_username
    DB_PASSWORD = module.rds.db_password
    
    # Cognito設定
    COGNITO_USER_POOL_ID     = module.cognito.user_pool_id
    COGNITO_CLIENT_ID        = module.cognito.user_pool_client_id
    COGNITO_USER_POOL_DOMAIN = module.cognito.user_pool_domain
    COGNITO_REGION           = var.aws_region
    
    # JWT設定
    JWT_ACCESS_SECRET  = "dev-access-secret-key-change-in-production"
    JWT_REFRESH_SECRET = "dev-refresh-secret-key-change-in-production"
    
    # 環境設定
    ENVIRONMENT = "development"
    AWS_REGION  = var.aws_region
  }
  sensitive = true
}
