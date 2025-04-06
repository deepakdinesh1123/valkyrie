variable "location" {
  type = string
}
# variable "security_groups" {
#   type = list(string)
# }
variable "spot_instance" {
  type = bool
}
variable "key_pair_name" {
  type = string
}
variable "access_point_mount_point" {
  type = string
}

variable "ebs_size" {
  type = number
}
variable "multi_attach_enabled" {
  type = bool
}
