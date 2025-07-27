# data "aws_vpc" "vpc" {
#   filter {
#     name   = "tag:Name"
#     values = ["valnix-vpc-useast-1"]
#   }
# }

# data "aws_subnet" "compute_subnet" {
#   filter {
#     name   = "tag:Name"
#     values = ["snet-compute-01"]
#   }
# }
data "aws_iam_role" "spot-fleet" {
  name = "AWSServiceRoleForEC2SpotFleet"
}