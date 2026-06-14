resource "aws_ecr_repository" "social_api" {
  name                 = var.project_name
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}

# Keep the repo from growing unbounded: drop untagged images (left behind
# when a tag like "latest" is reassigned to a newer build) after 7 days.
resource "aws_ecr_lifecycle_policy" "social_api" {
  repository = aws_ecr_repository.social_api.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Expire untagged images after 7 days"
        selection = {
          tagStatus   = "untagged"
          countType   = "sinceImagePushed"
          countUnit   = "days"
          countNumber = 7
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}
