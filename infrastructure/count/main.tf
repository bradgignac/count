variable "aws_region" {}

module "iam" {
  source = "../modules/iam"
}

module "dynamo" {
  source = "../modules/dynamo"
}

output "lambda_function_role_id" {
  value = "${module.iam.lambda_function_role_id}"
}
