resource "aws_db_subnet_group" "this" {
  name       = "${var.project_name}-db"
  subnet_ids = data.aws_subnets.default.ids
}

resource "aws_security_group" "rds" {
  name        = "${var.project_name}-rds"
  description = "Allow Postgres access from the ECS service only"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    description     = "Postgres from ECS service"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.ecs_service.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "random_password" "db" {
  length  = 32
  special = false
}

resource "aws_db_instance" "this" {
  identifier     = "${var.project_name}-db"
  engine         = "postgres"
  engine_version = var.db_engine_version
  instance_class = var.db_instance_class

  allocated_storage = var.db_allocated_storage
  storage_type      = "gp3"

  db_name  = var.db_name
  username = var.db_username
  password = random_password.db.result

  db_subnet_group_name   = aws_db_subnet_group.this.name
  vpc_security_group_ids = [aws_security_group.rds.id]

  publicly_accessible     = false
  multi_az                = false
  backup_retention_period = 1
  skip_final_snapshot     = true
}

# Encrypts the SSM parameter holding the DB connection string below.
resource "aws_kms_key" "ssm" {
  description             = "Encrypts ${var.project_name} SSM parameters"
  deletion_window_in_days = 7
}

locals {
  db_addr = "postgres://${var.db_username}:${random_password.db.result}@${aws_db_instance.this.address}:5432/${var.db_name}?sslmode=require"
}

resource "aws_ssm_parameter" "db_addr" {
  name   = "/${var.project_name}/db_addr"
  type   = "SecureString"
  key_id = aws_kms_key.ssm.key_id
  value  = local.db_addr
}

# Lets the ECS task execution role decrypt and read the DB connection string
# at task startup, and inject it as the DB_ADDR env var (see ecs.tf).
data "aws_iam_policy_document" "ecs_task_execution_secrets" {
  statement {
    actions   = ["ssm:GetParameters"]
    resources = [aws_ssm_parameter.db_addr.arn]
  }

  statement {
    actions   = ["kms:Decrypt"]
    resources = [aws_kms_key.ssm.arn]
  }
}

resource "aws_iam_role_policy" "ecs_task_execution_secrets" {
  name   = "${var.project_name}-ecs-task-execution-secrets"
  role   = aws_iam_role.ecs_task_execution.id
  policy = data.aws_iam_policy_document.ecs_task_execution_secrets.json
}
