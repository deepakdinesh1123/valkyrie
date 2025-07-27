output "ec2_public_ip" {
  value = var.deploy == true ? aws_instance.ec2[0].public_ip : 0
}
output "ec2_id" {
  value = var.deploy == true ? aws_instance.ec2[0].id : 0
}
output "ec2_availability_zone" {
  value = aws_instance.ec2[0].availability_zone
}