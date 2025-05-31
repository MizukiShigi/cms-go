# サービス情報
output "service_url" {
  description = "The URL of the Cloud Run service"
  value       = google_cloud_run_service.main.status[0].url
}

output "repository_url" {
  description = "The URL of the Artifact Registry repository"
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.main.repository_id}"
}

# GitHub Secrets用の必須値
output "workload_identity_provider" {
  description = "The Workload Identity Provider for GitHub Actions (WIF_PROVIDER)"
  value       = google_iam_workload_identity_pool_provider.github_provider.name
}

output "github_actions_service_account" {
  description = "The service account email for GitHub Actions (WIF_SERVICE_ACCOUNT)"
  value       = google_service_account.github_actions.email
}