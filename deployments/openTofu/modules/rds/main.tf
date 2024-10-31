resource "aws_db_subnet_group" "default" {
  name       = var.db_subnet_group_name
  subnet_ids = var.subnet_ids 
}

resource "aws_db_instance" "default" {
  identifier              = "mydbinstance"
  allocated_storage       = var.allocated_storage
  engine                 = var.rds_engine
  engine_version         = var.engine_version
  instance_class         = var.db_compute_instance
  db_subnet_group_name   = aws_db_subnet_group.default.name
  vpc_security_group_ids = var.security_group_ids
  username               = var.db_username
  password               = var.db_password
  db_name                = var.db_name
  skip_final_snapshot    = var.skip_create_final_snapshot
  multi_az               = var.multi_availability_zones
}
