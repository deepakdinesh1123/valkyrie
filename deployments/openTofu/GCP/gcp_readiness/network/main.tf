module "vpc" {
  source = "../../modules/vpc"

  vpc_name = "valkyrie-vpc"
}

module "subnet1" {
  source = "../../modules/vpc/subnet"

  snet_name = "valkyrie-worker-snet-01"
  snet_cidr = "10.0.0.0/22"
  region = "asia-south1"
  vpc_id = module.vpc.vpc_id
}

module "subnet2" {
  source = "../../modules/vpc/subnet"

  snet_name = "valkyrie-server-snet-01"
  snet_cidr = "10.0.4.0/22"
  region = "asia-south1"
  vpc_id = module.vpc.vpc_id
}

module "subnet3" {
  source = "../../modules/vpc/subnet"

  snet_name = "valkyrie-db-snet-01"
  snet_cidr = "10.0.8.0/22"
  region = "asia-south1"
  vpc_id = module.vpc.vpc_id
}

module "ssh_firewall" {
  source = "../../modules/vpc/firewall"

  vpc_id = module.vpc.vpc_id
  rule_name = "allow-ssh"
  ports = [ "22" ]
  protocol_name = "tcp"
  source_ips = ["0.0.0.0/0"]
}

module "db_firewall" {
  source = "../../modules/vpc/firewall"

  vpc_id = module.vpc.vpc_id
  rule_name = "allow-db"
  ports = [ "5432" ]
  protocol_name = "tcp"
  source_ips = ["0.0.0.0/0"]
}