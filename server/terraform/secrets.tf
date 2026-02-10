# Secrets Manager secret
resource "aws_secretsmanager_secret" "app_secrets" {
  name        = "${var.project_name}-${var.environment}-secrets"
  description = "Application secrets for ${var.project_name}"

  tags = {
    Name = "${var.project_name}-secrets"
  }
}

# Secret version with the actual values
resource "aws_secretsmanager_secret_version" "app_secrets" {
  secret_id = aws_secretsmanager_secret.app_secrets.id
  secret_string = jsonencode({
    MONGO_URI  = "mongodb+srv://christopherl4n_db_user:rHy9TfrKNPIkDK5m@cluster0.jb4hydn.mongodb.net/?appName=Cluster0"
    JWT_SECRET = "super_secret_change_me"
  })
}
