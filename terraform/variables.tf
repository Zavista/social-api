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

variable "db_name" {
  description = "Name of the application database"
  type        = string
  default     = "socialnetwork"
}

variable "db_username" {
  description = "Master username for the RDS instance"
  type        = string
  default     = "social_api_admin"
}

variable "db_instance_class" {
  description = "RDS instance size (db.t4g.micro is RDS free-tier eligible)"
  type        = string
  default     = "db.t4g.micro"
}

variable "db_allocated_storage" {
  description = "Allocated storage for RDS, in GB"
  type        = number
  default     = 20
}

variable "db_engine_version" {
  description = "Postgres engine version"
  type        = string
  default     = "16.4"
}
