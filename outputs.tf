output "lambda_function_name" {
  description = "部署的Lambda函数名称"
  value       = module.lambda_function.lambda_function_name
}

output "lambda_function_arn" {
  description = "Lambda函数的ARN"
  value       = module.lambda_function.lambda_function_arn
}

output "lambda_invoke_arn" {
  description = "Lambda函数的调用ARN"
  value       = module.lambda_function.lambda_function_invoke_arn
}

output "lambda_role_arn" {
  description = "Lambda函数使用的IAM角色ARN"
  value       = module.lambda_function.lambda_role_arn
}

output "lambda_role_name" {
  description = "Lambda函数使用的IAM角色名称"
  value       = module.lambda_function.lambda_role_name
}

output "lambda_invoke_command" {
  description = "调用Lambda函数的AWS CLI命令示例"
  value       = "aws lambda invoke --function-name ${module.lambda_function.lambda_function_name} --payload '{}' response.json --region ${var.aws_region} --endpoint-url https://lambda.${var.aws_region}.amazonaws.com.cn"
}