resource "github_actions_secret" "gcp_project_id" {
  repository      = "cms-go"
  secret_name     = "GCP_PROJECT_ID"
  plaintext_value = var.project_id
}

resource "github_actions_secret" "wif_provider" {
  repository      = "cms-go"
  secret_name     = "WIF_PROVIDER"
  plaintext_value = google_iam_workload_identity_pool_provider.github_provider.name
}

resource "github_actions_secret" "wif_service_account" {
  repository      = "cms-go"
  secret_name     = "WIF_SERVICE_ACCOUNT"
  plaintext_value = google_service_account.github_actions.email
}