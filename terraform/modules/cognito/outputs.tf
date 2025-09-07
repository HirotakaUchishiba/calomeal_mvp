# terraform/modules/cognito/outputs.tf

output "user_pool_id" {
  description = "ID of the Cognito User Pool"
  value       = aws_cognito_user_pool.main.id
}

output "user_pool_arn" {
  description = "ARN of the Cognito User Pool"
  value       = aws_cognito_user_pool.main.arn
}

output "user_pool_endpoint" {
  description = "Endpoint of the Cognito User Pool"
  value       = aws_cognito_user_pool.main.endpoint
}

output "user_pool_client_id" {
  description = "ID of the Cognito User Pool Client"
  value       = aws_cognito_user_pool_client.main.id
}

output "user_pool_client_secret" {
  description = "Secret of the Cognito User Pool Client"
  value       = aws_cognito_user_pool_client.main.client_secret
  sensitive   = true
}

output "user_pool_domain" {
  description = "Domain of the Cognito User Pool"
  value       = aws_cognito_user_pool_domain.main.domain
}

output "user_pool_domain_cloudfront_distribution_arn" {
  description = "CloudFront distribution ARN for the Cognito User Pool domain"
  value       = aws_cognito_user_pool_domain.main.cloudfront_distribution_arn
}

output "identity_pool_id" {
  description = "ID of the Cognito Identity Pool (if enabled)"
  value       = var.enable_identity_pool ? aws_cognito_identity_pool.main[0].id : null
}

output "identity_pool_arn" {
  description = "ARN of the Cognito Identity Pool (if enabled)"
  value       = var.enable_identity_pool ? aws_cognito_identity_pool.main[0].arn : null
}

# 環境変数用の出力
output "cognito_config" {
  description = "Cognito configuration for environment variables"
  value = {
    COGNITO_USER_POOL_ID     = aws_cognito_user_pool.main.id
    COGNITO_CLIENT_ID        = aws_cognito_user_pool_client.main.id
    COGNITO_USER_POOL_DOMAIN = aws_cognito_user_pool_domain.main.domain
    COGNITO_REGION           = data.aws_region.current.name
  }
}

# 現在のリージョン情報
data "aws_region" "current" {}
