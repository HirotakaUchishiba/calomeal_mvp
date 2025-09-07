# terraform/modules/cognito/variables.tf

variable "project_name" {
  description = "Name of the project"
  type        = string
  default     = "calomeal"
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
}

variable "callback_urls" {
  description = "List of callback URLs for the Cognito User Pool Client"
  type        = list(string)
  default     = ["http://localhost:5173", "https://localhost:5173"]
}

variable "logout_urls" {
  description = "List of logout URLs for the Cognito User Pool Client"
  type        = list(string)
  default     = ["http://localhost:5173", "https://localhost:5173"]
}

variable "enable_identity_pool" {
  description = "Whether to create a Cognito Identity Pool"
  type        = bool
  default     = false
}

variable "custom_domain" {
  description = "Custom domain for Cognito User Pool (optional)"
  type        = string
  default     = null
}

variable "email_from_name" {
  description = "From name for Cognito emails"
  type        = string
  default     = "CaloMeal"
}

variable "email_from_email" {
  description = "From email for Cognito emails"
  type        = string
  default     = "noreply@calomeal.com"
}

variable "sms_authentication_message" {
  description = "SMS authentication message template"
  type        = string
  default     = "Your authentication code is {####}"
}

variable "sms_verification_message" {
  description = "SMS verification message template"
  type        = string
  default     = "Your verification code is {####}"
}

variable "email_verification_message" {
  description = "Email verification message template"
  type        = string
  default     = "Please click the link below to verify your email address. {##Verify Email##}"
}

variable "email_verification_subject" {
  description = "Email verification subject"
  type        = string
  default     = "Verify your email address"
}

variable "tags" {
  description = "Additional tags to apply to resources"
  type        = map(string)
  default     = {}
}
