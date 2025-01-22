variable "mig_name" {
  
}
variable "vm_template_id" {
  
}
variable "mig_target_size" {
  
}
variable "zone_name" {
  type = string
  validation {
    condition = contains(["asia-south1-c", "asia-south1-b", "asia-south1-a"], var.zone_name)
    error_message = "Only asia-south1-c, asia-south1-b, asia-south1-a are allowed"
  }
}
variable "data_disk_name" {
  
}