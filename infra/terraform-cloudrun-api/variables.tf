variable "project_id" {
  type = string
  description = "The ID of the GCP project"
}

variable "region" {
  type = string
  description = "The region to deploy the Cloud Run API"
}

variable "image_name" {
  description = "Image name"
  type        = string
  default     = "cms-go"
}

variable "image_tag" {
  description = "Image tag"
  type        = string
  default     = "latest"
}

variable "db_name" {
  description = "Database name"
  type        = string
  default     = "cms"
}

variable "db_user" {
  description = "Database user name"
  type        = string
  default     = "postgres"
}

variable "db_tier" {
  description = "Database instance tier"
  type        = string
  default     = "db-f1-micro"
}

variable "db_activation_policy" {
  description = "Database activation policy (ALWAYS, NEVER)"
  type        = string
  default     = "ALWAYS"
}

variable "github_repository" {
  description = "GitHub repository in format 'owner/repo-name'"
  type        = string
}

variable "github_token" {
  description = "GitHub Personal Access Token"
  type        = string
  sensitive   = true
}