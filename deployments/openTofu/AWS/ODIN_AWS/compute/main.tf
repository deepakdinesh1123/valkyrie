module "vpc" {
  source = "../../modules/vpc"

  vpc_name          = "valnix-vpc-useast-1"
  vpc_address_block = "10.0.0.0/16"
}

module "compute_subnet" {
  source = "../../modules/vpc/subnet"

  subnet_name             = "snet-compute-01"
  vpc_id                  = module.vpc.vpc_id
  snet_availability_zone  = var.snet_availability_zone1
  snet_cidr               = "10.0.0.0/24"
  map_public_ip_on_launch = true
}

module "db_subnet01" {
  source = "../../modules/vpc/subnet"

  subnet_name             = "snet-db-01"
  vpc_id                  = module.vpc.vpc_id
  snet_availability_zone  = var.snet_availability_zone1
  snet_cidr               = "10.0.1.0/24"
  map_public_ip_on_launch = true
}

module "db_subnet02" {
  source = "../../modules/vpc/subnet"

  subnet_name             = "snet-db-02"
  vpc_id                  = module.vpc.vpc_id
  snet_availability_zone  = var.snet_availability_zone2
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

module "ssh_security_group" {
  source = "../../modules/vpc/securityGroups"

  security_grp_name = "SSH_Inbound"
  vpc_id            = module.vpc.vpc_id
  ingress_rules     = [
    {
      from_port   = 22
      to_port     = 22
      protocol    = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    },
    {

      from_port   = 2049
      to_port     = 2049
      protocol    = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    },
    {

      from_port   = 80
      to_port     = 80
      protocol    = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    }
    ,
    {

      from_port   = 443
      to_port     = 443
      protocol    = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    }
  ]
}

# resource "aws_network_interface" "my_network_interface" {
#   subnet_id   = module.subnet.subnet_id
#   private_ips = ["10.0.1.10"]

module "key_Pair" {
  source = "../../modules/keyPair"

  key_pair_name = var.key_pair_name
}

module "ec2" {
  source = "../../modules/ec2"
  deploy = true

  instance_type      = var.ec2_instance_type
  ami                = "ami-0866a3c8686eaeeba"
  security_group_ids = [module.ssh_security_group.sg_id]
  subnet_id          = module.compute_subnet.subnet_id
  #   nic_id = aws_network_interface.my_network_interface.id
  key_pair_name = module.key_Pair.aws_key_pair_name
  associate_pip = true
}

module "ebs" {
  source = "../../modules/ebs"

  ebs_name              = "shared_nix_store"
  ebs_size              = var.ebs_size
  ec2_availability_zone = module.ec2.ec2_availability_zone
  multi_attach_enabled  = var.multi_attach_enabled
  ebs_type              = var.ebs_type
  ebs_iops              = var.ebs_iops
  snapshot_id           = var.shared_nix_store_id #"snap-0daa1ed513204e62f"
}

module "ec2_spot_fleet" {
  source = "../../modules/ec2SpotFleet"

  aws_arn            = data.aws_iam_role.spot-fleet.arn
  ami_id             = "ami-0866a3c8686eaeeba"
  instance_types     = ["m6a.xlarge", "m5.xlarge", "m7a.xlarge", "m5a.xlarge", "m6i.xlarge", "m6id.xlarge"]
  key_Pair           = module.key_Pair.aws_key_pair_name
  subnet_id          = module.compute_subnet.subnet_id
  availability_zone  = module.ebs.ebs_availability_zone #data.aws_subnet.compute_subnet.availability_zone
  associate_pip      = true
  security_group_ids = [module.ssh_security_group.sg_id]

}

# module "server_ebs_vol_attach" {
#   source = "../../modules/ebs/ebs_volume_attach"

#   ec2_id     = module.ec2_spot_fleet.spot_fleet_id
#   volume_id  = module.ebs.ebs_id
#   depends_on = [module.ec2_spot_fleet]
# }


module "db_security_group" {
  source = "../../modules/vpc/securityGroups"

  security_grp_name = "DB_Inbound"
  vpc_id            = module.vpc.vpc_id
  ingress_rules     = [
    {
      from_port   = 22
      to_port     = 22
      protocol    = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    },
    {

      from_port   = 2049
      to_port     = 2049
      protocol    = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    },
    {

      from_port   = 5432
      to_port     = 5432
      protocol    = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    }
  ] 
}

module "random_password" {
  source             = "../../modules/randomPassword"
  deploy             = true
  pwd_include_lower  = true
  pwd_include_number = true
  pwd_include_upper  = true
  pwd_length         = 15
  pwd_special        = false
}
module "rds" {
  source = "../../modules/rds"

  db_instance_name           = "odinstagingserver"
  rds_engine                 = "postgres"
  engine_version             = "16.3"
  db_compute_instance        = var.rds_compute_type
  skip_create_final_snapshot = true
  db_name                    = "odinstagingdb"
  db_password                = module.random_password.password
  db_username                = "odinstagingowner"
  subnet_ids                 = [module.db_subnet01.subnet_id, module.db_subnet02.subnet_id]
  security_group_ids         = [module.db_security_group.sg_id]
  multi_availability_zones   = false
  allocated_storage          = 10

  db_subnet_group_name = "valnix-rds-subnet-group"
}

module "odin_secret" {
  source = "../../modules/secret"

  secret_name = "odin-worker-secrets-3"
  secret_value = {
    POSTGRES_USER = "odinstagingowner"
    POSTGRES_DB = "odinstagingdb"
    POSTGRES_PASSWORD = module.rds.rds_password
    POSTGRES_PORT = "5432"
    POSTGRES_HOST = local.rds_endpoint
  }

}
