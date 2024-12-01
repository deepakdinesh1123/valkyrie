data "google_service_account" "name" {
  account_id = "valk_user1"
}

data "google_compute_network" "vpc" {
  name = "valyrie-vpc"
}

# Fetch an existing Subnet
data "google_compute_subnetwork" "worker-subnet" {
  name   = "valkyrie-worker-snet-01"
  region = "asia-south1"
}

data "google_compute_subnetwork" "server-subnet" {
  name   = "valkyrie-server-snet-01"
  region = "asia-south1"
}

data "google_compute_subnetwork" "db-subnet" {
  name   = "valkyrie-db-snet-01"
  region = "asia-south1"
}