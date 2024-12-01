module "vpc" {
  source = "../../modules/vpc"

  vpc_name = "valnix-vpc-useast-1"
  vpc_address_block = "10.0.0.0/16"
}

module "compute_subnet" {
  source = "../../modules/vpc/subnet"

  subnet_name            = "snet-compute-01"
  vpc_id                 = module.vpc.vpc_id
  snet_availability_zone = "us-east-1a"
  snet_cidr              = "10.0.0.0/24"
  map_public_ip_on_launch = true
}

module "db_subnet01" {
  source = "../../modules/vpc/subnet"

  subnet_name             = "snet-db-01"
  vpc_id                  = module.vpc.vpc_id
  snet_availability_zone  = "us-east-1a"
  snet_cidr               = "10.0.1.0/24"
  map_public_ip_on_launch = true
}

module "db_subnet02" {
  source = "../../modules/vpc/subnet"

  subnet_name             = "snet-db-02"
  vpc_id                  = module.vpc.vpc_id
  snet_availability_zone  = "us-east-1b"
  snet_cidr               = "10.0.2.0/24"
  map_public_ip_on_launch = true
}

module "internet_gateway" {
  source = "../../modules/vpc/internetGateway"

  vpc_id = module.vpc.vpc_id
}

module "compute_route_table" {
  source = "../../modules/vpc/routeTable"

  route_table         = "compute-route-01"
  vpc_id              = module.vpc.vpc_id
  subnet_id           = module.compute_subnet.subnet_id
  internet_gateway_id = module.internet_gateway.internet_gateway_id
}

module "db_route_table01" {
  source = "../../modules/vpc/routeTable"

  route_table         = "db-route-01"
  vpc_id              = module.vpc.vpc_id
  subnet_id           = module.db_subnet01.subnet_id
  internet_gateway_id = module.internet_gateway.internet_gateway_id
}

module "db_route_table02" {
  source = "../../modules/vpc/routeTable"

  route_table         = "db-route-02"
  vpc_id              = module.vpc.vpc_id
  subnet_id           = module.db_subnet02.subnet_id
  internet_gateway_id = module.internet_gateway.internet_gateway_id
}