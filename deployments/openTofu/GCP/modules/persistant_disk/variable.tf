variable "disk_name" {
  
}
variable "disk_size_gb" {
  type = number
}
variable "project_name" {
  
}
variable "zone_name" {
  type = string
  validation {
    condition = contains(["asia-south1-c", "asia-south1-b", "asia-south1-a"], var.zone_name)
    error_message = "Only asia-south1-c, asia-south1-b, asia-south1-a are allowed"
  }
}
variable "disk_type" {
  
}
variable "disk_throughput" {
  type = number
}
variable "disk_iops" {
  type = number
}