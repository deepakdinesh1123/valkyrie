resource "google_compute_subnetwork" "subnet" {
  name          = var.snet_name
  ip_cidr_range = var.snet_cidr
  region        = var.region
  network       = var.vpc_id
}