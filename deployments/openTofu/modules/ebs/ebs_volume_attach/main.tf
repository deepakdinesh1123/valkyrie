resource "aws_volume_attachment" "ebs_att" {
  device_name = "/dev/sdh"
  volume_id   = var.volume_id
  instance_id = var.ec2_id
}
