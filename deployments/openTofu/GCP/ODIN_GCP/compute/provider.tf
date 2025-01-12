terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.0"
    }
  }

  required_version = ">= 1.5.0"

  # backend "gcs" {
  #   bucket  = "tofu-statefile-bucket"
  #   prefix  = "tofustate/compute"
  # }
}

provider "google" {
  credentials = "./keys.json"
  project     = var.project_name
  region      = var.location
}   