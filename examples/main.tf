# Copyright RoseSecurity 2026

terraform {
  required_version = ">= 1.5.6"
  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "~> 2.5"
    }
  }
}

# Uses only a couple of the settable arguments — filename + content.
resource "local_file" "example" {
  filename = "${path.module}/out.txt"
  content  = "hello"
}
