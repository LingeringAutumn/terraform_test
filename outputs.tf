output "lambda_function_name" {
  description = "部署的Lambda函数名称"
  value       = aws_lambda_function.function.function_name
}

output "lambda_function_arn" {
  description = "Lambda函数的ARN"
  value       = aws_lambda_function.function.arn
}

output "lambda_invoke_url" {
  description = "Lambda函数的基本调用URL"
  value       = "可通过AWS CLI使用: aws lambda invoke --function-name ${aws_lambda_function.function.function_name} --payload '{}' response.json"
}

output "lambda_role_arn" {
  description = "Lambda函数使用的IAM角色ARN"
  value       = aws_iam_role.lambda_role.arn
}