location                = "us-east-1"
snet_availability_zone1 = "us-east-1a"
snet_availability_zone2 = "us-east-1b"
key_pair_name           = "ec2_key_pair"

ebs_size             = 80
multi_attach_enabled = true
ebs_iops             = 1000
ebs_type             = "io1"

ec2_instance_type = "t3.micro"
rds_compute_type  = "db.t3.micro"