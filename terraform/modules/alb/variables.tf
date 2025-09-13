# terraform/modules/alb/variables.tf

variable "project_name" {
  description = "Name of the project"
  type        = string
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID where the ALB will be created"
  type        = string
}

variable "public_subnet_ids" {
  description = "List of public subnet IDs for the ALB"
  type        = list(string)
}

variable "target_port" {
  description = "Port on which targets receive traffic"
  type        = number
  default     = 8080
}

variable "ssl_certificate_arn" {
  description = "ARN of the SSL certificate for HTTPS listener"
  type        = string
  default     = null
}
