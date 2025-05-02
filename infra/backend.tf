terraform {
  backend "s3" {
    bucket         = "terraform-dev-state-bucket-229418028078"
    key            = "dev/infra/terraform.tfstate"
    region         = "eu-west-1"
    dynamodb_table = "terraform-dev-state-lock-229418028078"
    encrypt        = true
  }
}
