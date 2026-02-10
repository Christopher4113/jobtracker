# Build trigger to compile Go binary via CodeBuild
resource "null_resource" "build_trigger" {
  triggers = {
    always_run = timestamp()
  }

  provisioner "local-exec" {
    interpreter = ["/bin/bash", "-c"]
    command     = <<-EOF
      set -e

      echo "=== Starting build trigger ==="

      # Navigate to server directory
      cd ${path.module}/..

      # Create zip excluding terraform and hidden files
      echo "Creating source zip..."
      zip -r /tmp/server.zip . -x "terraform/*" -x ".*" -x "*.zip"

      # Upload to S3
      echo "Uploading source to S3..."
      aws s3 cp /tmp/server.zip s3://${aws_s3_bucket.artifacts.bucket}/source/server.zip

      # Start build and capture build ID
      echo "Starting CodeBuild..."
      BUILD_ID=$(aws codebuild start-build --project-name ${aws_codebuild_project.go_build.name} --query 'build.id' --output text --region us-east-1)
      echo "Build started with ID: $BUILD_ID"

      # Wait for build to complete (poll every 10 seconds, timeout after 10 minutes)
      TIMEOUT=600
      ELAPSED=0
      while true; do
        STATUS=$(aws codebuild batch-get-builds --ids $BUILD_ID --query 'builds[0].buildStatus' --output text --region us-east-1)

        if [ "$STATUS" = "SUCCEEDED" ]; then
          echo "Build succeeded!"
          break
        elif [ "$STATUS" = "FAILED" ] || [ "$STATUS" = "FAULT" ] || [ "$STATUS" = "STOPPED" ] || [ "$STATUS" = "TIMED_OUT" ]; then
          echo "Build failed with status: $STATUS"
          # Get build logs for debugging
          aws codebuild batch-get-builds --ids $BUILD_ID --query 'builds[0].phases[*].[phaseType,phaseStatus]' --output table --region us-east-1 || true
          rm -f /tmp/server.zip
          exit 1
        fi

        echo "Build status: $STATUS - waiting..."
        sleep 10
        ELAPSED=$((ELAPSED + 10))

        if [ $ELAPSED -ge $TIMEOUT ]; then
          echo "Build timed out after $TIMEOUT seconds"
          rm -f /tmp/server.zip
          exit 1
        fi
      done

      # Clean up
      rm -f /tmp/server.zip
      echo "=== Build trigger complete ==="
    EOF
  }

  depends_on = [
    aws_codebuild_project.go_build,
    aws_s3_bucket.artifacts
  ]
}

# Data source to reference the built Lambda artifact
data "aws_s3_object" "lambda_artifact" {
  bucket     = aws_s3_bucket.artifacts.bucket
  key        = "builds/function.zip"
  depends_on = [null_resource.build_trigger]
}
