resource "random_password" "rand_pwd" {
  count             = var.deploy ? 1 : 0
  length            = var.pwd_length
  override_special  = "@#!%"
  special           = var.pwd_special
  upper             = var.pwd_include_upper
  lower             = var.pwd_include_lower
  numeric           = var.pwd_include_number
}