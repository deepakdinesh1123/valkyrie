module "security_group" {
  source = "../../modules/vpc/securityGroups"

  security_grp_name = "SSH_Inbound"
  vpc_id            = data.aws_vpc.vpc.id
}

# resource "aws_network_interface" "my_network_interface" {
#   subnet_id   = module.subnet.subnet_id
#   private_ips = ["10.0.1.10"]

module "key_Pair" {
  source = "../../modules/keyPair"

  key_pair_name = var.key_pair_name
}

# module "ec2" {
#   source = "../../modules/ec2"
#   deploy = var.spot_instance == false ? false : true

#   instance_type      = "t3.micro"
#   ami                = "ami-0866a3c8686eaeeba"
#   security_group_ids = [module.security_group.sg_id]
#   subnet_id          = module.subnet.subnet_id
#   #   nic_id = aws_network_interface.my_network_interface.id
#   key_pair_name = module.key_Pair.aws_key_pair_name
#   associate_pip = true
# }

# module "ec2_spot" {
#   source = "../../modules/ec2SpotInstance"
#   deploy = var.spot_instance == true ? true : false

#   instance_type      = "t3.micro" # use only m5, c5, r5, t3, and z1d faily vms
#   ami                = "ami-0866a3c8686eaeeba"
#   spot_price         = "0.03"
#   spot_type          = "persistent"
#   security_group_ids = module.security_group.sg_id
#   subnet_id          = module.subnet.subnet_id
#   key_pair_name      = module.key_Pair.aws_key_pair_name
# }

# module "ebs" {
#   source = "../../modules/ebs"

#   ebs_size              = var.ebs_size
#   ec2_availability_zone = "us-east-1a" #module.ec2_spot.spot_ec2_availability_zone
#   multi_attach_enabled  = var.multi_attach_enabled
#   ebs_type              = "io1"
#   ebs_iops              = 1000
# }

module "ec2_spot_fleet" {
  source = "../../modules/ec2SpotFleet"

  ami_id                = "ami-010e773a908e799c1" 
  instance_types        = ["c5.large", "m5.large", "t3.large"]
  key_Pair              = module.key_Pair.aws_key_pair_name
  subnet_id             = data.aws_subnet.compute_subnet.id
  availability_zone     = data.aws_subnet.compute_subnet.availability_zone
  associate_pip         = true
  security_group_ids    = [ module.security_group.sg_id ]
}

# module "ebs_vol_attach" {
#   source = "../../modules/ebs/ebs_volume_attach"

#   ec2_id = var.spot_instance == true ? aws_spot_fleet_request.example.id : ""
#   volume_id = module.ebs.ebs_id
#   depends_on = [ aws_spot_fleet_request.example ]

# }

# module "efs" {
#   source = "../../modules/efs"

#   security_group_id = module.security_group.sg_id
#   subnet_id         = module.subnet.subnet_id
# }

# resource "null_resource" "configure_nfs" {
#   connection {
#     type        = "ssh"
#     user        = "ubuntu"
#     private_key = module.key_Pair.aws_key_pem
#     host        = var.spot_instance == true ? module.ec2_spot.spot_ec2_pip : module.ec2.ec2_public_ip
#   }
#   provisioner "remote-exec" {
#     inline = [

#       "sudo apt-get update -y",
#       "sudo mkdir -p /mnt/nixstore",
#       "sudo mount -t efs -o tls,accesspoint=${module.efs.access_point_id} ${module.efs.efs_id}:/ ${var.access_point_mount_point}"
#     ]
#   }
# }

# c7a.xlarge
# m5.large
# t3.large
# c5.xlarge