# CodeBuild project for building Go binary
resource "aws_codebuild_project" "go_build" {
  name          = "${var.project_name}-go-build"
  description   = "Build Go binary for Lambda deployment"
  build_timeout = 15
  service_role  = aws_iam_role.codebuild.arn

  artifacts {
    type = "NO_ARTIFACTS"
  }

  environment {
    compute_type                = "BUILD_GENERAL1_SMALL"
    image                       = "aws/codebuild/amazonlinux2-aarch64-standard:3.0"
    type                        = "ARM_CONTAINER"
    image_pull_credentials_type = "CODEBUILD"

    environment_variable {
      name  = "S3_BUCKET"
      value = aws_s3_bucket.artifacts.bucket
    }

    environment_variable {
      name  = "S3_KEY"
      value = "builds/function.zip"
    }
  }

  source {
    type      = "NO_SOURCE"
    buildspec = <<-EOF
      version: 0.2
      phases:
        install:
          runtime-versions:
            golang: 1.22
        pre_build:
          commands:
            - echo "Preparing build environment..."
            - mkdir -p /tmp/build
            - aws s3 cp s3://$S3_BUCKET/source/server.tar.gz /tmp/build/server.tar.gz
            - cd /tmp/build && tar -xzf server.tar.gz
        build:
          commands:
            - echo "Building Go binary..."
            - cd /tmp/build
            - go mod download
            - GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap main.go
        post_build:
          commands:
            - echo "Packaging and uploading..."
            - cd /tmp/build
            - zip function.zip bootstrap
            - aws s3 cp function.zip s3://$S3_BUCKET/$S3_KEY
            - echo "Build complete. Artifact uploaded to s3://$S3_BUCKET/$S3_KEY"
    EOF
  }

  tags = {
    Name = "${var.project_name}-go-build"
  }
}
