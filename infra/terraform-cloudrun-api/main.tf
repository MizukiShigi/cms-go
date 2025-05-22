resource "google_storage_bucket" "example" {
  name = "${var.project_id}-terraform-test-bucket"
  location = var.region
}


