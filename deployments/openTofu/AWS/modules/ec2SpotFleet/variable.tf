variable "key_Pair" {
  type = string
}
variable "subnet_id" {
  type = string
}
variable "availability_zone" {
  type = string
}
variable "instance_types" {
  type = list(string)
}
variable "associate_pip" {
  type = bool
}
variable "security_group_ids" {
  type = list(string)
}
variable "ami_id" {
  type = string
}