resource "aws_ebs_volume" "ebs" {
  availability_zone = var.ec2_availability_zone
  size              = var.ebs_size
  multi_attach_enabled = var.multi_attach_enabled
  type = var.ebs_type
  iops = var.ebs_iops
  snapshot_id = ""#var.snapshot_id

  tags = {
    Name = var.ebs_name
  }
}