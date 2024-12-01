resource "aws_subnet" "snet" {
  vpc_id            = var.vpc_id
  cidr_block        = var.snet_cidr
  availability_zone = var.snet_availability_zone
  map_public_ip_on_launch = var.map_public_ip_on_launch

  tags = {
    Name = var.subnet_name
  }
}