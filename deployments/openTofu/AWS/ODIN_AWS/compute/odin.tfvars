location                = "us-east-1"
snet_availability_zone1 = "us-east-1a"
snet_availability_zone2 = "us-east-1b"
key_pair_name           = "ec2_key_pair"

ebs_size             = 80
multi_attach_enabled = false
ebs_iops             = 5000
ebs_type             = "gp3"

ec2_instance_type = "t3.large"
rds_compute_type  = "db.t3.micro"