resource "google_sql_database_instance" "main" {
  name             = "${var.project_id}-postgres"
  database_version = "POSTGRES_17"
  region           = var.region
  
  deletion_protection = false

  settings {
    tier = var.db_tier
    activation_policy = var.db_activation_policy
    
    disk_size = 10
    disk_type = "PD_HDD"
    
    backup_configuration {
      enabled = false
    }
    
    ip_configuration {
      ipv4_enabled = true
      authorized_networks {
        value = "0.0.0.0/0"
        name  = "all"
      }
    }
  }
}

resource "google_sql_database" "database" {
  name     = var.db_name
  instance = google_sql_database_instance.main.name
}