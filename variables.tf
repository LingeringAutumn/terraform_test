variable "aws_region" {
  description = "部署AWS Lambda函数的区域"
  type        = string
  default     = "cn-northwest-1"  # 默认为东京区域，您可以根据需要修改
}

# 删除aws_access_key、aws_secret_key和aws_session_token变量

variable "lambda_function_name" {
  description = "Lambda函数的名称"
  type        = string
  default     = "terraform-lambda-demo"
}

variable "lambda_runtime" {
  description = "Lambda函数的运行时环境"
  type        = string
  default     = "nodejs18.x"
}

variable "lambda_timeout" {
  description = "Lambda函数的超时时间（秒）"
  type        = number
  default     = 30
}

variable "lambda_environment_variables" {
  description = "Lambda函数的环境变量"
  type        = map(string)
  default = {
    ENV = "dev"
  }
}