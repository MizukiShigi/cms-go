variable "project_id" {
  type = string
  description = "The ID of the GCP project"
}

variable "region" {
  type = string
  description = "The region to deploy the Cloud Run API"
}

variable "container_image" {
  description = "Container image URL"
  type        = string
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