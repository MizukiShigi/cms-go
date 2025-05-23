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

resource "google_cloud_run_service" "main" {
  name     = "${var.project_id}-api"
  location = var.region

  template {
    spec {
      containers {
        image = "gcr.io/cloudrun/hello"
        
        ports {
          container_port = 8080
        }

        # リソース制限
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