output "spot_ec2_pip" {
  value = var.deploy == true ? aws_spot_instance_request.vps[0].public_ip : ""
}
output "spot_ec2_availability_zone" {
  value = length(aws_spot_instance_request.vps) > 0 ? aws_spot_instance_request.vps[0].availability_zone : ""
}
output "spot_ec2_id" {
  value = length(aws_spot_instance_request.vps) > 0 ? aws_spot_instance_request.vps[0].spot_instance_id : ""
  description = "The ID of the Spot Instance."
}
