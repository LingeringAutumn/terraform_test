variable "aws_region" {
  description = "部署AWS Lambda函数的区域"
  type        = string
  default     = "ap-northeast-1"  # 默认为东京区域，您可以根据需要修改
}

variable "aws_access_key" {
  description = "AWS访问密钥ID"
  type        = string
  default     = ""  # 不要在代码中硬编码，建议使用环境变量或.tfvars文件
}

variable "aws_secret_key" {
  description = "AWS私有访问密钥"
  type        = string
  default     = ""  # 不要在代码中硬编码，建议使用环境变量或.tfvars文件
}

variable "aws_session_token" {
  description = "AWS会话令牌（使用临时凭证时需要）"
  type        = string
  default     = ""  # 对于大多数永久凭证，此值可以为空
}

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