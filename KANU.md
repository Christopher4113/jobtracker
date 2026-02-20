# KANU.md - Keep All Knowledge Updated

## Setup

- Install Terraform CLI (required for infrastructure deployment)
- Install Go 1.23+ (required for building the server)
- Run `go mod tidy` in `server/` directory to resolve dependencies after updating go.mod

## Deploy

- Navigate to `server/terraform/` directory
- Run `terraform init` to initialize Terraform backend and providers
- Run `terraform apply` to deploy AWS infrastructure (DynamoDB, Lambda, API Gateway, CodeBuild, S3, Secrets Manager)
