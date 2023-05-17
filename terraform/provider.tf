terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "4.47.0"
    }
  }
}

provider "google" {
  # Configuration options
  project     = var.gcp_project-id
  region      = var.gcp_region
}