# AWS配置
# 不要在此文件中存储敏感凭证

# 部署AWS Lambda函数的区域
aws_region = "cn-northwest-1"  # 默认为中国宁夏区域

# Lambda函数的名称
lambda_function_name = "feishu-notification-service"

# Lambda函数的运行时环境
# 使用provided.al2用于Go语言在Amazon Linux 2上运行
lambda_runtime = "provided.al2"

# Lambda函数的超时时间（秒）
lambda_timeout = 60

# 飞书机器人配置
feishu_app_id = "cli_a8028d701ebad00c"
feishu_app_secret = "PMHEgtwRpyQicbsnMbuCyhdeIpuat2cX"
feishu_base_url = "https://open.feishu.cn"

# 默认群聊ID（可选，部署后通过 /chats 接口获取实际群聊ID）
default_chat_id = ""

# 项目到群聊的映射（JSON格式字符串）
# 部署后先调用 /chats 接口获取群聊ID，然后更新此配置
project_chat_mapping = "{\"default\":\"\"}"

# Lambda函数的环境变量
lambda_environment_variables = {
  ENV = "prod"
  SERVICE_NAME = "feishu-notification-service"
  LOG_LEVEL = "info"
}
