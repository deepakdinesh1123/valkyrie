resource "aws_efs_access_point" "access-point" {
  file_system_id = var.efs_id

  posix_user {
    gid = 1000
    uid = 1000
  }

  root_directory {
    path = "/nix"
    creation_info {
      owner_gid   = 1000
      owner_uid   = 1000
      permissions = "0775"
    }
  }
}