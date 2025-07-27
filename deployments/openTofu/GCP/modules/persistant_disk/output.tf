output "disk_name" {
  value = google_compute_disk.persistant_disk.name
}
output "disk_id" {
  value = google_compute_disk.persistant_disk.id
}