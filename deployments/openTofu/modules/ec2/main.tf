resource "aws_instance" "ec2" {

  count                  = var.deploy == false ? 1 : 0
  ami                    = var.ami
  instance_type          = var.instance_type
  vpc_security_group_ids = var.security_group_ids
  subnet_id              = var.subnet_id
  # network_interface {
  #   delete_on_termination = true
  #   device_index = 0
  #   network_interface_id = var.nic_id
  # }
  key_name                    = var.key_pair_name
  associate_public_ip_address = var.associate_pip
}
