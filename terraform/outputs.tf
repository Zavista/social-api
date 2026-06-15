output "ecr_repository_url" {
  description = "Repository URL to docker push/pull against"
  value       = aws_ecr_repository.social_api.repository_url
}

output "alb_dns_name" {
  description = "Public URL of the API (via the ALB)"
  value       = "http://${aws_lb.this.dns_name}"
}

output "rds_endpoint" {
  description = "RDS connection endpoint (host:port) - for running migrations"
  value       = aws_db_instance.this.endpoint
}
