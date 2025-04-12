terraform {
  backend "s3" {
    bucket         = "valnix-terraform-state-bucket"
    key            = "AWS_Readiness/network/network01.tfstate"
    region         = "us-east-1"
    dynamodb_table = "tf-backend"
    encrypt        = true
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
    time = {
      source  = "hashicorp/time"
      version = "~> 0.9"
    }
  }

  required_version = ">= 1.0"
}

provider "aws" {
  region     = "us-east-1"
  access_key = ""
  secret_key = ""
  # assume_role {
  #   role_arn     = "arn:aws:iam::775188627313:user/manoj_ha"
  # }
}
