
output "lambda_function_role_id" {
  value = "${aws_iam_role.count_function.arn}"
}
