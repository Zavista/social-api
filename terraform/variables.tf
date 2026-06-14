variable "aws_region" {
  description = "AWS region to deploy into"
  type        = string
  default     = "us-east-1"
}

variable "project_name" {
  description = "Name used for naming/tagging resources"
  type        = string
  default     = "social-api"
}
