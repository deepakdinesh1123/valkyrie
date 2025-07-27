resource "google_compute_disk" "persistant_disk" {
  name       = var.disk_name
  size       = var.disk_size_gb
  type       = var.disk_type
  zone       = var.zone_name
  physical_block_size_bytes = 4096
  provisioned_iops           = contains(["hyperdisk"],var.disk_type) ? var.disk_iops : 0
  provisioned_throughput     = contains(["hyperdisk"],var.disk_type) ? var.disk_throughput : 0
}
