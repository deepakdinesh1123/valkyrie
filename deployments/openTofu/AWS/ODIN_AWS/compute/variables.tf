variable "location" {
  type = string
}
variable "snet_availability_zone1" {
  type = string
}
variable "snet_availability_zone2" {
  type = string
}
# variable "security_groups" {
#   type = list(string)
# }
variable "key_pair_name" {
  type = string
}
variable "ebs_size" {
  type = number
  default = 80
}
variable "multi_attach_enabled" {
  type = bool
}
variable "ebs_type" {
  type = string
}
variable "ebs_iops" {
  type = number
}
variable "shared_nix_store_id" {
  type    = string
  default = ""
}
variable "ec2_instance_type" {
  type = string
}
variable "rds_compute_type" {
  type = string
}