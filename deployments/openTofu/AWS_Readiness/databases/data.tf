data "aws_vpc" "vpc" {
  filter {
    name   = "tag:Name"
    values = ["valnix-vpc-useast-1"]
  }
}

data "aws_subnet" "db_subnet01" {
  filter {
    name   = "tag:Name"
    values = ["snet-db-01"]
  }
}

data "aws_subnet" "db_subnet02" {
  filter {
    name   = "tag:Name"
    values = ["snet-db-02"]
  }
}
