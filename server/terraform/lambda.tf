# Lambda function for Go Fiber server
resource "aws_lambda_function" "api" {
  function_name = "${var.project_name}-api"
  description   = "Job Tracker API - Go Fiber server"

  # Use S3 artifact from CodeBuild
  s3_bucket         = aws_s3_bucket.artifacts.bucket
  s3_key            = "builds/function.zip"
  s3_object_version = data.aws_s3_object.lambda_artifact.version_id

  # Runtime configuration for Go
  runtime       = "provided.al2023"
  architectures = ["arm64"]
  handler       = "bootstrap"

  # Performance settings
  memory_size = 256
  timeout     = 30

  # Attach Lambda execution role
  role = aws_iam_role.lambda_execution.arn

  # AWS Lambda Web Adapter layer (us-east-1 ARM64)
  layers = ["arn:aws:lambda:us-east-1:753240598075:layer:LambdaAdapterLayerArm64:24"]

  environment {
    variables = {
      # Required for Lambda Web Adapter
      AWS_LAMBDA_EXEC_WRAPPER = "/opt/bootstrap"
      PORT                    = "8080"

      # Application config
      CORS_ORIGIN = "*"

      # Secrets ARN for SDK-based retrieval
      SECRETS_ARN = aws_secretsmanager_secret.app_secrets.arn

      # Pass secrets directly as environment variables
      MONGO_URI  = jsondecode(aws_secretsmanager_secret_version.app_secrets.secret_string)["MONGO_URI"]
      JWT_SECRET = jsondecode(aws_secretsmanager_secret_version.app_secrets.secret_string)["JWT_SECRET"]
    }
  }

  depends_on = [
    null_resource.build_trigger,
    aws_iam_role_policy.lambda_logs,
    aws_iam_role_policy.lambda_secrets
  ]

  tags = {
    Name = "${var.project_name}-api"
  }
}

# CloudWatch Log Group for Lambda
resource "aws_cloudwatch_log_group" "lambda" {
  name              = "/aws/lambda/${var.project_name}-api"
  retention_in_days = 14

  tags = {
    Name = "${var.project_name}-api-logs"
  }
}
