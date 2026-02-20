# KANU.md - Keep All Knowledge Updated

## Setup

- Install Terraform CLI (required for infrastructure deployment)

## Deploy

- Navigate to `server/terraform/` directory
- Run `terraform init` to initialize Terraform backend and providers
- Run `terraform apply` to deploy AWS infrastructure (DynamoDB, Lambda, API Gateway, CodeBuild, S3, Secrets Manager)
