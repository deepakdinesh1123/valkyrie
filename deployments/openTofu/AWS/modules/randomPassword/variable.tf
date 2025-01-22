variable "deploy" {
  type    = bool
}

variable "pwd_length" {
  type = number
  default = 26
}

variable "pwd_special" {
  type = bool
  default = true
}

variable "pwd_include_upper" {
  type = bool
  default = true
}

variable "pwd_include_lower" {
  type = bool
  default = true
}

variable "pwd_include_number" {
  type = bool
  default = true
}
