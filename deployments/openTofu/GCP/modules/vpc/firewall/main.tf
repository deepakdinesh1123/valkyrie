resource "google_compute_firewall" "allow_ssh" {
  name    = var.rule_name
  network = var.vpc_id

  allow {
    protocol = var.protocol_name
    ports    = var.ports
  }

  source_ranges = var.source_ips
}
