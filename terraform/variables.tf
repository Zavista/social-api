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

variable "container_port" {
  description = "Port the API container listens on"
  type        = number
  default     = 8080
}

variable "image_tag" {
  description = "Docker image tag to deploy (must already exist in ECR, see `make ecr-push`)"
  type        = string
  default     = "0.0.1"
}

variable "task_cpu" {
  description = "Fargate task vCPU units (256 = 0.25 vCPU)"
  type        = number
  default     = 256
}

variable "task_memory" {
  description = "Fargate task memory, in MB"
  type        = number
  default     = 512
}

variable "desired_count" {
  description = "Number of ECS tasks to run"
  type        = number
  default     = 1
}
