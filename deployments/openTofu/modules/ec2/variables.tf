variable "instance_type" {
  type = string
}
variable "ami" {
  type = string
}
variable "security_group_ids" {
  type = list(string)
}
variable "subnet_id" {
  type = string
}
variable "associate_pip" {
  type = bool
}
variable "deploy" {
  type = bool
}
variable "key_pair_name" {
  type = string
}