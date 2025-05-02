provider "aws" {
  region = var.aws_region
}

resource "aws_ecr_repository" "go_api_repo" {
  name = "go-api"

  image_tag_mutability = "MUTABLE"
  image_scanning_configuration {
    scan_on_push = false
  }
}