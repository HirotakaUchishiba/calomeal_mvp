# terraform/modules/cognito/main.tf

# Cognito User Pool
resource "aws_cognito_user_pool" "main" {
  name = "${var.project_name}-user-pool"

  # パスワードポリシー
  password_policy {
    minimum_length    = 8
    require_lowercase = true
    require_uppercase = true
    require_numbers   = true
    require_symbols   = true
  }

  # ユーザー属性
  schema {
    name                = "email"
    attribute_data_type = "String"
    required            = true
    mutable             = true
  }

  schema {
    name                = "name"
    attribute_data_type = "String"
    required            = false
    mutable             = true
  }

  # 認証設定
  username_attributes = ["email"]
  auto_verified_attributes = ["email"]

  # アカウント回復設定
  account_recovery_setting {
    recovery_mechanism {
      name     = "verified_email"
      priority = 1
    }
  }

  # メール設定
  email_configuration {
    email_sending_account = "COGNITO_DEFAULT"
  }

  # 管理者作成ユーザー設定
  admin_create_user_config {
    allow_admin_create_user_only = false
    invite_message_template {
      email_message = "Welcome to ${var.project_name}! Your username is {username} and temporary password is {####}."
      email_subject = "Welcome to ${var.project_name}"
      sms_message   = "Your username is {username} and temporary password is {####}."
    }
  }

  # デバイス設定
  device_configuration {
    challenge_required_on_new_device      = true
    device_only_remembered_on_user_prompt = false
  }

  # ユーザープールのタグ
  tags = {
    Name        = "${var.project_name}-user-pool"
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

  # リフレッシュトークン設定
  refresh_token_validity {
    refresh_token_validity = 7
  }

  # 認証設定
  generate_secret                      = false  # SPA用にfalse
  prevent_user_existence_errors        = "ENABLED"
  enable_token_revocation              = true
  enable_propagate_additional_user_context_data = false

  # コールバックURL設定
  callback_urls = var.callback_urls
  logout_urls   = var.logout_urls

  # サポートされるOAuthフロー
  supported_identity_providers = ["COGNITO"]

  # OAuth設定
  allowed_oauth_flows = ["code", "implicit"]
  allowed_oauth_scopes = [
    "email",
    "openid",
    "profile"
  ]
  allowed_oauth_flows_user_pool_client = true

  # セキュリティ設定
  token_validity_units {
    access_token  = "hours"
    id_token      = "hours"
    refresh_token = "days"
  }

  # ユーザープールクライアントのタグ
  tags = {
    Name        = "${var.project_name}-user-pool-client"
    Environment = var.environment
    Project     = var.project_name
  }
}

# Cognito User Pool Domain
resource "aws_cognito_user_pool_domain" "main" {
  domain       = "${var.project_name}-${var.environment}-auth"
  user_pool_id = aws_cognito_user_pool.main.id
}

# Cognito Identity Pool (オプション - 将来的な拡張用)
resource "aws_cognito_identity_pool" "main" {
  count = var.enable_identity_pool ? 1 : 0

  identity_pool_name               = "${var.project_name}-identity-pool"
  allow_unauthenticated_identities = false

  cognito_identity_providers {
    client_id               = aws_cognito_user_pool_client.main.id
    provider_name           = aws_cognito_user_pool.main.endpoint
    server_side_token_check = true
  }

  tags = {
    Name        = "${var.project_name}-identity-pool"
    Environment = var.environment
    Project     = var.project_name
  }
}
