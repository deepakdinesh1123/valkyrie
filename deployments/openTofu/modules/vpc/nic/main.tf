resource "aws_network_interface" "my_network_interface" {
  subnet_id   = var.subnet_id
  private_ips = ["10.0.1.10"]  # Specify a private IP if needed
}
