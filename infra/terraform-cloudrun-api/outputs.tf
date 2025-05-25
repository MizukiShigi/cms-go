output "project_id" {
  description = "The ID of the GCP project"
  value       = var.project_id
}

output "region" {
  description = "The region of the GCP project"
  value       = var.region
}

output "service_url" {
  description = "The URL of the Cloud Run service"
  value       = google_cloud_run_service.main.status[0].url
}

output "service_name" {
  description = "The name of the Cloud Run service"
  value       = google_cloud_run_service.main.name
}

output "repository_url" {
  description = "The URL of the Artifact Registry repository"
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.main.repository_id}"
}