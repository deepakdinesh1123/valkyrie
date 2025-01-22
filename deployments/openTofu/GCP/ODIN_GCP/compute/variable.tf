variable "location" {
  type = string
}
variable "availability_zone" {
  type = string
}
variable "project_name" {
  type = string
}
variable "disk_iops" {
  type = number
}
variable "disk_size_gb" {
  type = number
}
variable "disk_throughput" {
  type = number
}
variable "worker_machine_type" {
  type = string
}
variable "server_machine_type" {
  type = string
}