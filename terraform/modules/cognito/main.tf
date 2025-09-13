# terraform/modules/cognito/main.tf

# Cognito User Pool
resource "aws_cognito_user_pool" "main" {
  name = "${var.project_name}-user-pool"

  # パスワードポリシー
  password_policy {
    minimum_length    = 8
    require_lowercase = true
    require_numbers   = true
    require_symbols   = true
    require_uppercase = true
  }

  # ユーザープールの設定
  username_attributes = ["email"]
  auto_verified_attributes = ["email"]

  # アカウント回復設定
  account_recovery_setting {
    recovery_mechanism {
      name     = "verified_email"
      priority = 1
    }
  }

  tags = {
    Environment = var.environment
    Project     = var.project_name
  }
}

# Cognito User Pool Client
resource "aws_cognito_user_pool_client" "main" {
  name         = "${var.project_name}-user-pool-client"
  user_pool_id = aws_cognito_user_pool.main.id

  # 認証フロー設定
  explicit_auth_flows = [
    "ALLOW_USER_PASSWORD_AUTH",
    "ALLOW_USER_SRP_AUTH",
    "ALLOW_REFRESH_TOKEN_AUTH"
  ]

  # トークン有効期限設定
  access_token_validity  = 1  # 1時間
  id_token_validity      = 1  # 1時間
  refresh_token_validity = 7  # 7日

  # コールバックURL設定
  callback_urls = var.callback_urls
  logout_urls   = var.logout_urls

  # OAuth設定
  allowed_oauth_flows = ["code", "implicit"]
  allowed_oauth_scopes = [
    "email",
    "openid",
    "profile"
  ]
  allowed_oauth_flows_user_pool_client = true

  # セキュリティ設定
  generate_secret = false  # SPA用にfalse
  prevent_user_existence_errors = "ENABLED"
  enable_token_revocation = true
}

# Cognito User Pool Domain
resource "aws_cognito_user_pool_domain" "main" {
  count        = var.custom_domain != null ? 1 : 0
  domain       = var.custom_domain
  user_pool_id = aws_cognito_user_pool.main.id
}

# デフォルトドメイン（カスタムドメインが設定されていない場合）
resource "aws_cognito_user_pool_domain" "default" {
  count        = var.custom_domain == null ? 1 : 0
  domain       = "${var.project_name}-${var.environment}-${random_string.domain_suffix.result}"
  user_pool_id = aws_cognito_user_pool.main.id
}

# ランダム文字列（デフォルトドメイン用）
resource "random_string" "domain_suffix" {
  length  = 8
  special = false
  upper   = false
}
