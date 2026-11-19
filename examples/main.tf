# Copyright RoseSecurity 2026

terraform {
  required_version = ">= 1.5.6"
  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "~> 2.5"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
  }
}

# Uses only a couple of the settable arguments — filename + content.
resource "local_file" "example" {
  filename = "${path.module}/out.txt"
  content  = "hello"
}

resource "aws_rds_cluster" "example" {
  cluster_identifier      = "example"
  engine                  = "aurora-mysql"
  engine_version          = "5.7.mysql_aurora.2.03.2"
  availability_zones      = ["us-west-2a", "us-west-2b", "us-west-2c"]
  database_name           = "exampledb"
  master_username         = "exampleuser"
  master_password         = "examplepassword"
  backup_retention_period = 5
  preferred_backup_window = "07:00-09:00"
}

