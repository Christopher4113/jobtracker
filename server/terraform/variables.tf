variable "project_name" {
  description = "Name of the project, used for resource naming"
  type        = string
  default     = "jobtracker"
}

variable "environment" {
  description = "Deployment environment (prod, dev, staging)"
  type        = string
  default     = "prod"
}
