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
  description   = "使用Terraform模块部署的Lambda函数"
  handler       = "bootstrap"  # Go使用bootstrap作为handler
  runtime       = var.lambda_runtime
  timeout       = var.lambda_timeout

  source_path = "${path.module}/src"

  environment_variables = var.lambda_environment_variables

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
    Project     = "terraform-lambda-demo"
  }
}

# 为任何需要的其他资源添加配置