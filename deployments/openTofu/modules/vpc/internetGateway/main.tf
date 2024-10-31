resource "aws_internet_gateway" "vpc-int-gw" {
  vpc_id = var.vpc_id
}