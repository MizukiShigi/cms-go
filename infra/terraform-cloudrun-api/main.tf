resource "google_storage_bucket" "example" {
  name = "${var.project_id}-terraform-test-bucket"
  location = var.region
}

resource "google_artifact_registry_repository" "main" {
  location      = var.region
  repository_id = "${var.project_id}-repo"
  description   = "Docker repository for API"
  format        = "DOCKER"
}

resource "google_service_account" "cloud_run" {
  account_id   = "cloud-run-service"
  display_name = "Cloud Run Service Account"
  description  = "Service account for Cloud Run API service"
}

# Secret Managerアクセス権限
resource "google_project_iam_member" "secret_accessor" {
  project = var.project_id
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:${google_service_account.cloud_run.email}"
}

# Cloud SQLクライアント権限
resource "google_project_iam_member" "cloudsql_client" {
  project = var.project_id
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${google_service_account.cloud_run.email}"
}

# ログ書き込み権限
resource "google_project_iam_member" "log_writer" {
  project = var.project_id
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${google_service_account.cloud_run.email}"
}

locals {
  container_image = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.main.repository_id}/${var.image_name}:${var.image_tag}"
}

resource "google_cloud_run_service" "main" {
  name     = "${var.project_id}-api"
  location = var.region

  template {
    spec {
      service_account_name = google_service_account.cloud_run.email
      containers {
        image = local.container_image
        
        ports {
          container_port = 8080
        }

        env {
          name  = "DB_HOST"
          value = google_sql_database_instance.main.public_ip_address
        }
        
        env {
          name  = "DB_NAME"
          value = google_sql_database.database.name
        }
        
        env {
          name  = "DB_USER"
          value = var.db_user
        }
        
        env {
          name = "DB_PASSWORD"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.db_password.secret_id
              key  = "latest"
            }
          }
        }
        
        env {
          name = "JWT_SECRET_KEY"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.jwt_secret.secret_id
              key  = "latest"
            }
          }
        }

        resources {
          limits = {
            cpu    = "1000m"  # 1 CPU
            memory = "512Mi"  # 512MB RAM
          }
        }

        env {
          name  = "ENV"
          value = "development"
        }
      }
    }

    metadata {
      annotations = {
        # 最大インスタンス数
        "autoscaling.knative.dev/maxScale" = "5"
        # CPUスロットリングを無効（パフォーマンス向上）
        "run.googleapis.com/cpu-throttling" = "false"
      }
    }
  }

  lifecycle {
    ignore_changes = [
      template[0].spec[0].containers[0].image
    ]
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_cloud_run_service_iam_member" "public" {
  service  = google_cloud_run_service.main.name
  location = google_cloud_run_service.main.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}

# === GitHub Actions Deployment Resources ===

# GitHub Actions用サービスアカウント
resource "google_service_account" "github_actions" {
  account_id   = "github-actions"
  display_name = "GitHub Actions Service Account"
  description  = "Service account for GitHub Actions CI/CD pipeline"
}

# Workload Identity Pool
resource "google_iam_workload_identity_pool" "github_pool" {
  workload_identity_pool_id = "github-pool"
  display_name             = "GitHub Actions Pool"
  description              = "Workload Identity Pool for GitHub Actions"
}

# Workload Identity Provider
resource "google_iam_workload_identity_pool_provider" "github_provider" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.github_pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "github-provider"
  display_name                       = "GitHub Actions Provider"
  description                        = "OIDC identity pool provider for GitHub Actions"

  attribute_mapping = {
    "google.subject"       = "assertion.sub"
    "attribute.actor"      = "assertion.actor"
    "attribute.repository" = "assertion.repository"
  }

  attribute_condition = "attribute.repository == '${var.github_repository}'"

  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }
}

# GitHub Actions必須権限
resource "google_project_iam_member" "github_actions_artifactregistry" {
  project = var.project_id
  role    = "roles/artifactregistry.writer"
  member  = "serviceAccount:${google_service_account.github_actions.email}"
}

resource "google_project_iam_member" "github_actions_cloudrun" {
  project = var.project_id
  role    = "roles/run.developer"
  member  = "serviceAccount:${google_service_account.github_actions.email}"
}

resource "google_project_iam_member" "github_actions_serviceaccount" {
  project = var.project_id
  role    = "roles/iam.serviceAccountUser"
  member  = "serviceAccount:${google_service_account.github_actions.email}"
}

# Workload Identity User binding
resource "google_service_account_iam_member" "github_actions_workload_identity" {
  service_account_id = google_service_account.github_actions.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/projects/${data.google_project.project.number}/locations/global/workloadIdentityPools/${google_iam_workload_identity_pool.github_pool.workload_identity_pool_id}/attribute.repository/${var.github_repository}"
}

# プロジェクト情報取得用
data "google_project" "project" {
  project_id = var.project_id
}