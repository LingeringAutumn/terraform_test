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

# API Gateway 输出
output "api_gateway_id" {
  description = "API Gateway的ID"
  value       = aws_api_gateway_rest_api.feishu_api.id
}

output "api_gateway_execution_arn" {
  description = "API Gateway的执行ARN"
  value       = aws_api_gateway_rest_api.feishu_api.execution_arn
}

output "api_gateway_url" {
  description = "API Gateway的调用URL"
  value       = "https://${aws_api_gateway_rest_api.feishu_api.id}.execute-api.${var.aws_region}.amazonaws.com.cn/${aws_api_gateway_stage.feishu_api_stage.stage_name}"
}

output "notification_endpoint" {
  description = "飞书通知接口端点URL"
  value       = "https://${aws_api_gateway_rest_api.feishu_api.id}.execute-api.${var.aws_region}.amazonaws.com.cn/${aws_api_gateway_stage.feishu_api_stage.stage_name}"
}

output "health_check_endpoint" {
  description = "健康检查接口端点URL"
  value       = "https://${aws_api_gateway_rest_api.feishu_api.id}.execute-api.${var.aws_region}.amazonaws.com.cn/${aws_api_gateway_stage.feishu_api_stage.stage_name}/health"
}

output "chats_list_endpoint" {
  description = "获取群聊列表接口端点URL"
  value       = "https://${aws_api_gateway_rest_api.feishu_api.id}.execute-api.${var.aws_region}.amazonaws.com.cn/${aws_api_gateway_stage.feishu_api_stage.stage_name}/chats"
}

# 测试命令示例
output "curl_health_check" {
  description = "健康检查cURL命令示例"
  value       = "curl -X GET \"https://${aws_api_gateway_rest_api.feishu_api.id}.execute-api.${var.aws_region}.amazonaws.com.cn/${aws_api_gateway_stage.feishu_api_stage.stage_name}/health\""
}

output "curl_notification_example" {
  description = "发送通知cURL命令示例"
  value = <<-EOT
curl -X POST "https://${aws_api_gateway_rest_api.feishu_api.id}.execute-api.${var.aws_region}.amazonaws.com.cn/${aws_api_gateway_stage.feishu_api_stage.stage_name}" \
  -H "Content-Type: application/json" \
  -d '{
    "message_type": "error",
    "title": "测试错误通知",
    "content": "这是一个测试错误消息",
    "source": {
      "service_name": "test-service",
      "environment": "dev"
    },
    "details": {
      "level": "ERROR",
      "timestamp": "'$(date -Iseconds)'"
    },
    "target": {
      "project_name": "project1"
    }
  }'
EOT
}

output "lambda_invoke_command" {
  description = "调用Lambda函数的AWS CLI命令示例"
  value       = "aws lambda invoke --function-name ${module.lambda_function.lambda_function_name} --payload '{}' response.json --region ${var.aws_region} --endpoint-url https://lambda.${var.aws_region}.amazonaws.com.cn"
}
