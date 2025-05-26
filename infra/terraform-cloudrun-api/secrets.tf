resource "google_secret_manager_secret" "db_password" {
  secret_id = "db-password"
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret" "jwt_secret" {
  secret_id = "jwt-secret-key"  
  replication {
    auto {}
  }
}