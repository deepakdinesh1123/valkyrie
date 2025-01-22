output "public_ip_address" {
  value = google_compute_instance.compute-engine.network_interface[0].access_config[0].nat_ip
}