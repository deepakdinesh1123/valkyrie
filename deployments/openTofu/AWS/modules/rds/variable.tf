variable "db_instance_name" {
  type = string
}

variable "rds_engine" {
  type = string
}

variable "subnet_ids" {
  type = list(string)
}

variable "db_username" {
  type = string
}

variable "db_password" {
  type = string
}

variable "db_name" {
  type = string
}

variable "skip_create_final_snapshot" {
  type = bool
}

variable "allocated_storage" {
  type = number
}

variable "security_group_ids" {
  type = list(string)
}

variable "multi_availability_zones" {
  type = bool
}

variable "db_compute_instance" {
  
}

variable "db_subnet_group_name" {
  
}

variable "engine_version" {
  type = string
}