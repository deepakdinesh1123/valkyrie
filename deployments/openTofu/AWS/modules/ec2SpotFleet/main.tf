resource "aws_spot_fleet_request" "example" {
  iam_fleet_role = var.aws_arn
  target_capacity = 1
	terminate_instances_on_delete = true

	dynamic "launch_specification" {
    for_each = var.instance_types
    content {
      ami                  = var.ami_id
      instance_type        = launch_specification.value
      key_name             = var.key_Pair
      spot_price           = "0.08"
      iam_instance_profile = "OdinBuilderWorker"
      subnet_id            = var.subnet_id
      availability_zone    = var.availability_zone
			associate_public_ip_address = var.associate_pip
			vpc_security_group_ids = var.security_group_ids

			user_data = <<-EOF
				#!/bin/bash

				cd /home/ubuntu
				echo "hello" > log.txt
				BUCKET_NAME="valnix-stage-bucket"
				REGION="us-east-1"
				SCRIPT_NAME="deploy.sh"
				SCRIPT_PATH="/home/ubuntu/$SCRIPT_NAME"

				wget -P /home/ubuntu/ https://valnix-stage-bucket.s3.amazonaws.com/deploy.sh

				cd /home/ubuntu

				chmod +x $SCRIPT_NAME

				su ubuntu -c ./$SCRIPT_NAME
			EOF
			root_block_device {
				volume_size           = "30"
				volume_type           = "gp3"
				delete_on_termination = "true"
			}
			tags = {
				Name = "valnix-spot-fleet-instance"
			}
		}
	}

	tags = {
			Name = "valnix-worker"
	}
}

# resource "aws_launch_template" "ec2_spot_fleet_template" {
#   image_id = "ami-0866a3c8686eaeeba"
#   name = "ec2-spot-fleet-template"
#   key_name = module.key_Pair.aws_key_pair_name
#   instance_requirements {
#     allowed_instance_types = [ 
#       "c5.xlarge",
#       "c6i.2xlarge",
#       "t3.large",
#       "m5.large",     # Added a general-purpose instance type
#       "m5a.large",
#       "c5.large",
#       "c6i.large"
#     ]

#     memory_mib {
#       min = 4096  # Lowered from 8192
#       max = 65536 # Increased to allow more options
#     }

#     vcpu_count {
#       min = 2     # Lowered from 4
#       max = 32    # Increased to allow larger instance types
#     }
#   }

#   network_interfaces {
#     subnet_id = module.subnet.subnet_id
#     associate_public_ip_address = true
#     security_groups = [ module.security_group.sg_id ]
#   }

#   placement {
#     availability_zone = module.ebs.ebs_availability_zone
#   }
# }

# resource "aws_spot_fleet_request" "ec2_spot_fleet" {
#   iam_fleet_role = "arn:aws:iam::775188627313:role/aws-service-role/spotfleet.amazonaws.com/AWSServiceRoleForEC2SpotFleet"
#   target_capacity = 1
#   spot_price = "0.6"

#   launch_template_config {
#     launch_template_specification {
#       id      = aws_launch_template.ec2_spot_fleet_template.id
#       version = aws_launch_template.ec2_spot_fleet_template.latest_version
#     }
#   }
#   tags = {
#     Name = "valnix-spot-fleet"
#   }
# }