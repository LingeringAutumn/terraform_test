provider "aws" {
  region = var.aws_region
  
  # 显式设置中国区域的终端节点URL
  endpoints {
    sts = "https://sts.${var.aws_region}.amazonaws.com.cn"
  }
}

# 使用AWS提供的Lambda模块
module "lambda_function" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "~> 4.0"

  function_name = var.lambda_function_name
  description   = "飞书通知服务Lambda函数"
  handler       = "bootstrap"  # Go使用bootstrap作为handler
  runtime       = var.lambda_runtime
  timeout       = var.lambda_timeout

  source_path = "${path.module}/src"

  # 合并环境变量
  environment_variables = merge(var.lambda_environment_variables, {
    FEISHU_APP_ID         = var.feishu_app_id
    FEISHU_APP_SECRET     = var.feishu_app_secret
    FEISHU_BASE_URL       = var.feishu_base_url
    DEFAULT_CHAT_ID       = var.default_chat_id
    PROJECT_CHAT_MAPPING  = var.project_chat_mapping
  })

  # 针对中国区域进行调整
  create_role = true
  lambda_at_edge = false
  
  # 为中国区域自定义IAM策略
  assume_role_policy_statements = {
    lambda = {
      effect  = "Allow",
      actions = ["sts:AssumeRole"],
      principals = {
        service = {
          type        = "Service",
          identifiers = ["lambda.amazonaws.com"]
        }
      }
    }
  }

  # 其他配置
  cloudwatch_logs_retention_in_days = 14
  attach_cloudwatch_logs_policy     = true
  
  tags = {
    Environment = lookup(var.lambda_environment_variables, "ENV", "dev")
    Project     = "feishu-notification-service"
    Service     = "notification"
  }
}

# API Gateway 配置
resource "aws_api_gateway_rest_api" "feishu_api" {
  name        = "${var.lambda_function_name}-api"
  description = "飞书通知服务API Gateway"

  endpoint_configuration {
    types = ["REGIONAL"]
  }

  tags = {
    Environment = lookup(var.lambda_environment_variables, "ENV", "dev")
    Project     = "feishu-notification-service"
  }
}

# API Gateway Resource - 健康检查
resource "aws_api_gateway_resource" "health" {
  rest_api_id = aws_api_gateway_rest_api.feishu_api.id
  parent_id   = aws_api_gateway_rest_api.feishu_api.root_resource_id
  path_part   = "health"
}

# API Gateway Resource - 群聊列表
resource "aws_api_gateway_resource" "chats" {
  rest_api_id = aws_api_gateway_rest_api.feishu_api.id
  parent_id   = aws_api_gateway_rest_api.feishu_api.root_resource_id
  path_part   = "chats"
}

# API Gateway Method - GET /health
resource "aws_api_gateway_method" "health_get" {
  rest_api_id   = aws_api_gateway_rest_api.feishu_api.id
  resource_id   = aws_api_gateway_resource.health.id
  http_method   = "GET"
  authorization = "NONE"
}

# API Gateway Method - GET /chats
resource "aws_api_gateway_method" "chats_get" {
  rest_api_id   = aws_api_gateway_rest_api.feishu_api.id
  resource_id   = aws_api_gateway_resource.chats.id
  http_method   = "GET"
  authorization = "NONE"
}

# API Gateway Method - POST /
resource "aws_api_gateway_method" "notify_post" {
  rest_api_id   = aws_api_gateway_rest_api.feishu_api.id
  resource_id   = aws_api_gateway_rest_api.feishu_api.root_resource_id
  http_method   = "POST"
  authorization = "NONE"
}

# API Gateway Integration - 健康检查
resource "aws_api_gateway_integration" "health_integration" {
  rest_api_id = aws_api_gateway_rest_api.feishu_api.id
  resource_id = aws_api_gateway_resource.health.id
  http_method = aws_api_gateway_method.health_get.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = module.lambda_function.lambda_function_invoke_arn
}

# API Gateway Integration - 群聊列表
resource "aws_api_gateway_integration" "chats_integration" {
  rest_api_id = aws_api_gateway_rest_api.feishu_api.id
  resource_id = aws_api_gateway_resource.chats.id
  http_method = aws_api_gateway_method.chats_get.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = module.lambda_function.lambda_function_invoke_arn
}

# API Gateway Integration - 通知接口
resource "aws_api_gateway_integration" "notify_integration" {
  rest_api_id = aws_api_gateway_rest_api.feishu_api.id
  resource_id = aws_api_gateway_rest_api.feishu_api.root_resource_id
  http_method = aws_api_gateway_method.notify_post.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = module.lambda_function.lambda_function_invoke_arn
}

# Lambda Permission for API Gateway
resource "aws_lambda_permission" "api_gw" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = module.lambda_function.lambda_function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_api_gateway_rest_api.feishu_api.execution_arn}/*/*"
}

# API Gateway Deployment
resource "aws_api_gateway_deployment" "feishu_api_deployment" {
  depends_on = [
    aws_api_gateway_method.health_get,
    aws_api_gateway_method.chats_get,
    aws_api_gateway_method.notify_post,
    aws_api_gateway_integration.health_integration,
    aws_api_gateway_integration.chats_integration,
    aws_api_gateway_integration.notify_integration,
  ]

  rest_api_id = aws_api_gateway_rest_api.feishu_api.id

  triggers = {
    redeployment = sha1(jsonencode([
      aws_api_gateway_resource.health.id,
      aws_api_gateway_resource.chats.id,
      aws_api_gateway_method.health_get.id,
      aws_api_gateway_method.chats_get.id,
      aws_api_gateway_method.notify_post.id,
      aws_api_gateway_integration.health_integration.id,
      aws_api_gateway_integration.chats_integration.id,
      aws_api_gateway_integration.notify_integration.id,
    ]))
  }

  lifecycle {
    create_before_destroy = true
  }
}

# API Gateway Stage
resource "aws_api_gateway_stage" "feishu_api_stage" {
  deployment_id = aws_api_gateway_deployment.feishu_api_deployment.id
  rest_api_id   = aws_api_gateway_rest_api.feishu_api.id
  stage_name    = lookup(var.lambda_environment_variables, "ENV", "dev")

  tags = {
    Environment = lookup(var.lambda_environment_variables, "ENV", "dev")
    Project     = "feishu-notification-service"
  }
}

# 为任何需要的其他资源添加配置
