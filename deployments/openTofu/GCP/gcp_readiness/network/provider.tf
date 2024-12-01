terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.0"
    }
  }

  required_version = ">= 1.5.0"

  backend "gcs" {
    bucket  = "tofu-statefile-bucket"
    prefix  = "tofustate/network"
    
  }
}

provider "google" {
  credentials = "./keys.json"
  project     = "valkyrie-project"
  region      = "asia-south1"
  zone        = "asia-south1-c"
}   