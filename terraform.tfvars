# AWS配置
# 不要在此文件中存储敏感凭证

# 部署AWS Lambda函数的区域
aws_region = "cn-northwest-1"  # 默认为中国宁夏区域

# Lambda函数的名称
# 默认值为"terraform-lambda-demo"，可以根据需要修改
lambda_function_name = "terraform-lambda"

# Lambda函数的运行时环境
# 使用provided.al2用于Go语言在Amazon Linux 2上运行
lambda_runtime = "provided.al2"

# Lambda函数的超时时间（秒）
# 默认值为30秒，可以根据函数的复杂度调整
lambda_timeout = 30

# Lambda函数的环境变量
# 默认值为{"ENV": "dev"}，可以根据需要添加更多键值对
lambda_environment_variables = {
  ENV = "dev"
}

# 其他配置（可选）
