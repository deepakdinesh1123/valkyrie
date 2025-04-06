resource "aws_route_table" "route-table-vps-env" {
  vpc_id = var.vpc_id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = var.internet_gateway_id
  }

  tags = {
    Name = var.route_table
  }
}

resource "aws_route_table_association" "subnet-association" {
  subnet_id      = var.subnet_id
  route_table_id = aws_route_table.route-table-vps-env.id
}
