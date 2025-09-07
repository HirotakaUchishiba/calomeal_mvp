variable "db_password" {
  description = "Database password"
  type        = string
  sensitive   = true
}

variable "cognito_custom_domain" {
  description = "Custom domain for Cognito User Pool (optional)"
  type        = string
  default     = null
}

variable "jwt_access_secret" {
  description = "JWT access token secret"
  type        = string
  sensitive   = true
}

variable "jwt_refresh_secret" {
  description = "JWT refresh token secret"
  type        = string
  sensitive   = true
}
