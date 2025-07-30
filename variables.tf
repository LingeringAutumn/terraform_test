variable "aws_region" {
  description = "部署AWS Lambda函数的区域"
  type        = string
  default     = "cn-northwest-1"  # 默认为中国宁夏区域
}

variable "lambda_function_name" {
  description = "Lambda函数的名称"
  type        = string
}

variable "lambda_runtime" {
  description = "Lambda函数的运行时环境"
  type        = string
  default     = "provided.al2"
}

variable "lambda_timeout" {
  description = "Lambda函数的超时时间（秒）"
  type        = number
  default     = 60
}

# 飞书机器人配置变量
variable "feishu_app_id" {
  description = "飞书机器人应用ID"
  type        = string
  sensitive   = true
}

variable "feishu_app_secret" {
  description = "飞书机器人应用密钥"
  type        = string
  sensitive   = true
}

variable "feishu_base_url" {
  description = "飞书API基础URL"
  type        = string
  default     = "https://open.feishu.cn"
}

variable "default_chat_id" {
  description = "默认群聊ID"
  type        = string
  default     = ""
}

variable "project_chat_mapping" {
  description = "项目名到群聊ID的映射（JSON格式字符串）"
  type        = string
  default     = "{\"default\":\"\"}"
}

variable "lambda_environment_variables" {
  description = "Lambda函数的环境变量"
  type        = map(string)
  default = {
    ENV = "dev"
  }
}
