# 飞书通知服务测试指南

## 快速回答您的问题

### ✅ 1. 飞书凭证已配置
您的飞书凭证已写入 `terraform.tfvars` 文件：
- App ID: `cli_a8028d701ebad00c`
- App Secret: `PMHEgtwRpyQicbsnMbuCyhdeIpuat2cX`

### ✅ 2. 自动获取群聊并选择发送
是的，这个Lambda函数支持：
- **自动获取群聊列表**: `GET /chats` 接口
- **灵活选择发送目标**: 支持三种方式选择群聊
  - 直接指定群聊ID: `"target": {"chat_ids": ["oc_xxx", "oc_yyy"]}`
  - 项目名映射: `"target": {"project_name": "user-service"}`
  - 使用默认群聊: 不指定target时使用默认群聊

### 3. 真实发送测试方法
**完整测试流程如下：**

### 4. 接口参数格式
**详细的请求/响应格式见下方**

---

## 完整部署和测试流程

### 第一步：部署Lambda函数

```bash
# 1. 进入项目目录
cd aws-lambda-terraform

# 2. 初始化Terraform
terraform init

# 3. 查看部署计划
terraform plan

# 4. 执行部署
terraform apply
# 输入 yes 确认部署

# 5. 记录部署输出的URL
terraform output
```

部署成功后，您会看到类似输出：
```
api_gateway_url = "https://abc123def.execute-api.cn-northwest-1.amazonaws.com.cn/prod"
chats_list_endpoint = "https://abc123def.execute-api.cn-northwest-1.amazonaws.com.cn/prod/chats"
health_check_endpoint = "https://abc123def.execute-api.cn-northwest-1.amazonaws.com.cn/prod/health"
notification_endpoint = "https://abc123def.execute-api.cn-northwest-1.amazonaws.com.cn/prod"
```

### 第二步：测试服务是否正常运行

```bash
# 1. 健康检查 - 确认Lambda函数正常
curl -X GET "https://your-api-gateway-url/prod/health"

# 预期响应：
# {
#   "status": "healthy",
#   "service": "feishu-notification-service", 
#   "timestamp": "2024-01-30T11:35:00Z"
# }
```

### 第三步：获取群聊列表

```bash
# 2. 获取群聊列表 - 确认飞书API连接正常
curl -X GET "https://your-api-gateway-url/prod/chats"

# 预期响应：
# {
#   "success": true,
#   "message": "Successfully retrieved 3 chats",
#   "timestamp": "2024-01-30T11:35:00Z",
#   "count": 3,
#   "chats": [
#     {
#       "chat_id": "oc_a1b2c3d4e5f6",
#       "name": "开发团队群",
#       "description": "",
#       "members": 15
#     },
#     {
#       "chat_id": "oc_f6e5d4c3b2a1", 
#       "name": "运维告警群",
#       "description": "",
#       "members": 8
#     }
#   ]
# }
```

### 第四步：更新群聊配置

根据第三步获取的群聊ID，更新 `terraform.tfvars` 文件：

```hcl
# 使用真实的群聊ID更新配置
project_chat_mapping = "{\"user-service\":\"oc_a1b2c3d4e5f6\",\"monitoring\":\"oc_f6e5d4c3b2a1\",\"default\":\"oc_a1b2c3d4e5f6\"}"
default_chat_id = "oc_a1b2c3d4e5f6"
```

然后重新部署：
```bash
terraform apply
```

### 第五步：测试真实消息发送

**测试1：发送到指定群聊ID**
```bash
curl -X POST "https://your-api-gateway-url/prod" \
  -H "Content-Type: application/json" \
  -d '{
    "message_type": "info",
    "title": "测试消息",
    "content": "这是一条测试消息，用于验证飞书通知服务是否正常工作",
    "target": {
      "chat_ids": ["oc_a1b2c3d4e5f6"]
    }
  }'
```

**测试2：通过项目名发送**
```bash
curl -X POST "https://your-api-gateway-url/prod" \
  -H "Content-Type: application/json" \
  -d '{
    "message_type": "error", 
    "title": "系统错误",
    "content": "用户服务出现数据库连接错误",
    "source": {
      "service_name": "user-service",
      "environment": "prod"
    },
    "details": {
      "level": "ERROR",
      "timestamp": "2024-01-30T11:35:00Z",
      "error_code": "DB_CONNECTION_FAILED"
    },
    "target": {
      "project_name": "user-service"
    }
  }'
```

**成功响应示例：**
```json
{
  "success": true,
  "message_id": "1706618100123456789",
  "message": "Sent to 1 chats, failed 0 chats",
  "timestamp": "2024-01-30T11:35:00Z",
  "details": {
    "processed_at": "2024-01-30T11:35:00Z",
    "sent_to_chats": ["oc_a1b2c3d4e5f6"],
    "failed_chats": [],
    "retry_scheduled": false
  }
}
```

---

## 完整接口参数格式

### 1. 通知接口 `POST /`

#### 请求格式 (最完整版本)
```json
{
  // === 基础信息 (必需) ===
  "message_type": "error",           // 必需：error/warning/info/success/alert
  "title": "系统错误告警",            // 可选：消息标题
  "content": "数据库连接失败，请立即处理", // 必需：消息内容
  
  // === 来源信息 (可选) ===
  "source": {
    "service_name": "user-service",   // 服务名称
    "module_name": "database",       // 模块名称
    "environment": "production",     // 环境：dev/test/staging/prod
    "region": "cn-northwest-1",      // 区域信息
    "version": "v1.2.3"             // 版本信息
  },
  
  // === 详细信息 (可选) ===
  "details": {
    "level": "ERROR",               // 级别：DEBUG/INFO/WARN/ERROR/FATAL
    "timestamp": "2024-01-30T11:35:00Z", // ISO8601时间戳
    "trace_id": "trace-abc123",     // 追踪ID
    "user_id": "user-456",          // 用户ID
    "request_id": "req-789",        // 请求ID
    "error_code": "DB_TIMEOUT",     // 错误代码
    "stack_trace": "at db.connect(db.go:45)\nat main.init(main.go:23)", // 堆栈信息
    "metadata": {                   // 额外元数据
      "connection_pool_size": 10,
      "retry_count": 3
    }
  },
  
  // === 目标配置 (必需其中一种) ===
  "target": {
    // 方式1：直接指定群聊ID列表
    "chat_ids": ["oc_a1b2c3d4e5f6", "oc_f6e5d4c3b2a1"],
    
    // 方式2：通过项目名映射 (与chat_ids二选一)
    "project_name": "user-service",
    
    // 可选配置
    "priority": "high",             // 优先级：low/normal/high/urgent
    "at_all": false,               // 是否@所有人
    "at_users": ["user1", "user2"] // @指定用户ID列表
  },
  
  // === 扩展字段 (为未来功能预留) ===
  "extensions": {
    "retry_count": 3,              // 重试次数
    "callback_url": "https://api.example.com/webhook", // 回调URL
    "tags": ["urgent", "database"], // 标签
    "custom_fields": {             // 自定义字段
      "ticket_id": "TICKET-123",
      "severity": "high"
    },
    "rate_limit_key": "user-service", // 限流key
    "deduplication_key": "error-db-connection-20240130" // 去重key
  }
}
```

#### 最简请求格式
```json
{
  "message_type": "error",
  "content": "这是一条错误消息",
  "target": {
    "project_name": "user-service"
  }
}
```

#### 响应格式
```json
{
  "success": true,                    // 是否成功
  "message_id": "1706618100123456789", // 消息ID
  "message": "Sent to 2 chats, failed 0 chats", // 状态描述
  "timestamp": "2024-01-30T11:35:00Z", // 处理时间
  "details": {
    "processed_at": "2024-01-30T11:35:00Z", // 处理时间
    "sent_to_chats": ["oc_a1b2c3d4e5f6"],   // 发送成功的群聊ID
    "failed_chats": [],                      // 发送失败的群聊ID
    "retry_scheduled": false                 // 是否安排重试
  }
}
```

### 2. 获取群聊列表接口 `GET /chats`

#### 请求
```bash
GET /chats
# 无需参数
```

#### 响应格式
```json
{
  "success": true,
  "message": "Successfully retrieved 3 chats",
  "timestamp": "2024-01-30T11:35:00Z",
  "count": 3,
  "chats": [
    {
      "chat_id": "oc_a1b2c3d4e5f6",     // 群聊ID (用于发送消息)
      "name": "开发团队群",               // 群聊名称
      "description": "开发团队日常沟通",   // 群聊描述
      "avatar": "https://...",           // 群聊头像URL
      "owner_id": "ou_owner123",         // 群主ID
      "members": 15                      // 成员数量
    }
  ]
}
```

### 3. 健康检查接口 `GET /health`

#### 请求
```bash
GET /health
# 无需参数
```

#### 响应格式
```json
{
  "status": "healthy",
  "service": "feishu-notification-service",
  "timestamp": "2024-01-30T11:35:00Z"
}
```

---

## 飞书消息格式示例

当您发送通知时，Lambda函数会将JSON格式化为易读的飞书消息：

### 错误消息格式
```
🚨 错误报警
标题: 系统错误告警
服务: user-service/database
环境: production(cn-northwest-1)
级别: ERROR
时间: 2024-01-30 19:35:00
内容: 数据库连接失败，请立即处理
追踪ID: trace-abc123
请求ID: req-789
错误码: DB_TIMEOUT
版本: v1.2.3
```

### 成功消息格式
```
✅ 成功通知
标题: 部署完成
服务: user-service
环境: production
时间: 2024-01-30 19:35:00
内容: 用户服务v1.2.3部署成功
版本: v1.2.3
```

---

## 故障排除

### 1. 如果健康检查失败
```bash
# 检查Lambda函数是否部署成功
terraform output lambda_function_name

# 查看Lambda日志
aws logs tail /aws/lambda/feishu-notification-service --follow --region cn-northwest-1
```

### 2. 如果获取群聊列表失败
- 检查飞书App ID和App Secret是否正确
- 确认机器人已被添加到群聊中
- 检查飞书应用权限是否包含 `im:chat` 和 `im:message`

### 3. 如果消息发送失败
- 确认使用的群聊ID是否正确 (通过 `/chats` 接口获取)
- 检查机器人是否在目标群聊中
- 查看CloudWatch日志了解详细错误信息

### 4. 查看详细日志
```bash
# 实时查看Lambda日志
aws logs tail /aws/lambda/feishu-notification-service --follow --region cn-northwest-1

# 查看API Gateway日志 (如果启用)
aws logs tail /aws/apigateway/feishu-notification-service-api --follow --region cn-northwest-1
```

---

## 高级使用技巧

### 1. 批量发送到多个群聊
```json
{
  "message_type": "alert",
  "content": "系统维护通知：将于今晚22:00-24:00进行维护",
  "target": {
    "chat_ids": ["oc_dev_team", "oc_ops_team", "oc_product_team"]
  }
}
```

### 2. 集成到CI/CD流水线
```bash
# Jenkins Pipeline示例
curl -X POST "${FEISHU_LAMBDA_URL}" \
  -H "Content-Type: application/json" \
  -d "{
    \"message_type\": \"success\",
    \"title\": \"部署成功\",
    \"content\": \"${JOB_NAME} 构建 #${BUILD_NUMBER} 部署成功\",
    \"source\": {
      \"service_name\": \"${JOB_NAME}\",
      \"environment\": \"${DEPLOY_ENV}\",
      \"version\": \"${BUILD_NUMBER}\"
    },
    \"target\": {
      \"project_name\": \"ci-cd\"
    }
  }"
```

### 3. 集成到应用程序
```python
import requests
import json
from datetime import datetime

def send_feishu_notification(message_type, content, service_name, error_details=None):
    payload = {
        "message_type": message_type,
        "content": content,
        "source": {
            "service_name": service_name,
            "environment": "prod"
        },
        "details": {
            "level": "ERROR" if message_type == "error" else "INFO",
            "timestamp": datetime.now().isoformat()
        },
        "target": {
            "project_name": service_name
        }
    }
    
    if error_details:
        payload["details"].update(error_details)
    
    response = requests.post(
        "https://your-api-gateway-url/prod",
        headers={"Content-Type": "application/json"},
        data=json.dumps(payload)
    )
    
    return response.json()

# 使用示例
send_feishu_notification(
    message_type="error",
    content="用户注册失败",
    service_name="user-service",
    error_details={
        "error_code": "VALIDATION_FAILED",
        "trace_id": "trace-123"
    }
)
```

这个Lambda飞书推送服务现在已经完全实现了您要求的功能：完全解耦、通用接口、高扩展性，并且可以真实发送消息到飞书群聊中。
