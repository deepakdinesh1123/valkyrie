locals {
  rds_endpoint = split(":", module.rds.rds_pip)[0]
}