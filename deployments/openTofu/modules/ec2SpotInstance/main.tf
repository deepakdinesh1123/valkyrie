resource "aws_spot_instance_request" "vps" {
  count                  = var.deploy == "true" ? 1 : 0
  ami                    = var.ami
  spot_price             = var.spot_price
  instance_type          = var.instance_type
  spot_type              = var.spot_type
  # block_duration_minutes = 120
  wait_for_fulfillment   = "true"
  key_name               = var.key_pair_name

  security_groups = [ var.security_group_ids ]
  subnet_id = var.subnet_id
  associate_public_ip_address = true
  instance_interruption_behavior = "terminate" 
}