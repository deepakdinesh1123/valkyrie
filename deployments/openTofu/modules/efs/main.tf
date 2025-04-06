resource "aws_efs_file_system" "efs" {}

resource "aws_efs_mount_target" "mount" {
  file_system_id  = aws_efs_file_system.efs.id
  security_groups = [var.security_group_id]
  subnet_id       = var.subnet_id
}

module "efs_access_point" {
  source = "./accessPoint"

  efs_id = aws_efs_file_system.efs.id
}
