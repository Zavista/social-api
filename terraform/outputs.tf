output "ecr_repository_url" {
  description = "Repository URL to docker push/pull against"
  value       = aws_ecr_repository.social_api.repository_url
}
