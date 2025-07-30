# 飞书通知服务 Lambda 函数

这是一个部署在AWS Lambda上的通用飞书通知服务，设计为一个解耦的"通知中心接口"。它可以接收来自不同模块的HTTP POST请求，并将这些通知格式化后发送到指定的飞书群聊中。

## 功能特性

### 🚀 核心功能
- **HTTP接口**: 提供RESTful API接口，支持通过HTTP POST请求发送通知
- **多消息类型**: 支持错误、警告、信息、成功、告警等多种消息类型
- **智能路由**: 根据项目名称或直接指定群聊ID来路由消息
- **格式化消息**: 自动格式化消息内容，包含丰富的上下文信息
- **健康检查**: 提供健康检查接口用于监控服务状态

### 🔧 技术特性
- **AWS Lambda**: 无服务器架构，按需执行，成本效益高
- **API Gateway**: 提供公网HTTP接口，支持HTTPS
- **Go语言**: 高性能、低延迟的处理能力
- **Terraform**: 基础设施即代码，自动化部署和管理
- **环境变量配置**: 灵活的配置管理，支持多环境部署

### 🌟 扩展性设计
- **通用接口**: 标准化的JSON接口，方便任何模块接入
- **预留字段**: 为未来功能预留了大量扩展字段
- **多群聊支持**: 支持同时发送到多个群聊
- **项目映射**: 灵活的项目到群聊映射配置

## 接口设计

### 1. 通知接口
**端点**: `POST /`  
**描述**: 发送通知消息到飞书群聊

#### 请求体结构
```json
{
  "message_type": "error",           // 必需：消息类型（error/warning/info/success/alert）
  "title": "错误标题",                // 可选：消息标题
  "content": "错误详细内容",           // 必需：消息内容
  
  // 来源信息
  "source": {
    "service_name": "user-service",   // 服务名称
    "module_name": "auth",           // 模块名称
    "environment": "prod",           // 环境（dev/test/staging/prod）
    "region": "cn-northwest-1",      // 区域
    "version": "v1.2.3"             // 版本
  },
  
  // 详细信息
  "details": {
    "level": "ERROR",                // 级别（DEBUG/INFO/WARN/ERROR/FATAL）
    "timestamp": "2024-01-15T10:30:00Z", // ISO8601时间戳
    "trace_id": "abc123",           // 追踪ID
    "user_id": "user123",           // 用户ID
    "request_id": "req456",         // 请求ID
    "error_code": "AUTH_001",       // 错误代码
    "stack_trace": "...",           // 堆栈信息
    "metadata": {}                  // 额外元数据
  },
  
  // 目标配置
  "target": {
    "chat_ids": ["oc_xxx", "oc_yyy"], // 直接指定群聊ID列表
    "project_name": "user-service",    // 或通过项目名映射
    "priority": "high",               // 优先级（low/normal/high/urgent）
    "at_all": false,                 // 是否@所有人
    "at_users": ["user1", "user2"]   // @指定用户
  },
  
  // 扩展字段（未来功能）
  "extensions": {
    "retry_count": 0,               // 重试次数
    "callback_url": "",             // 回调URL
    "tags": ["urgent", "db"],       // 标签
    "custom_fields": {},            // 自定义字段
    "rate_limit_key": "",           // 限流key
    "deduplication_key": ""         // 去重key
  }
}
```

#### 最简请求示例
```json
{
  "message_type": "error",
  "content": "数据库连接失败",
  "target": {
    "project_name": "user-service"
  }
}
```

#### 响应结构
```json
{
  "success": true,
  "message_id": "1674654600123456789",
  "message": "Sent to 2 chats, failed 0 chats",
  "timestamp": "2024-01-15T10:30:00Z",
  "details": {
    "processed_at": "2024-01-15T10:30:00Z",
    "sent_to_chats": ["oc_xxx", "oc_yyy"],
    "failed_chats": [],
    "retry_scheduled": false
  }
}
```

### 2. 健康检查接口
**端点**: `GET /health`  
**描述**: 检查服务健康状态

#### 响应示例
```json
{
  "status": "healthy",
  "service": "feishu-notification-service", 
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## 部署指南

### 前提条件

1. **安装工具**
   ```bash
   # 安装Terraform
   # 参考: https://terraform.io/downloads

   # 安装AWS CLI
   # 参考: https://aws.amazon.com/cli/
   ```

2. **AWS凭证配置**
   ```bash
   aws configure
   # 或设置环境变量
   export AWS_ACCESS_KEY_ID="your-access-key"
   export AWS_SECRET_ACCESS_KEY="your-secret-key"
   export AWS_DEFAULT_REGION="cn-northwest-1"
   ```

3. **飞书应用配置**
   - 在飞书开发者后台创建企业自建应用
   - 获取App ID和App Secret
   - 配置应用权限：
     - `im:chat` - 获取与更新群信息
     - `im:message` - 获取与发送单聊、群聊消息
   - 将机器人添加到需要接收通知的群聊中

### 配置步骤

1. **克隆代码**
   ```bash
   git clone <your-repo>
   cd aws-lambda-terraform
   ```

2. **配置变量**
   编辑 `terraform.tfvars` 文件：
   ```hcl
   # 基本配置
   lambda_function_name = "feishu-notification-service"
   aws_region = "cn-northwest-1"
   
   # 飞书配置
   feishu_app_id = "cli_xxxxxxxxxx"
   feishu_app_secret = "your-secret-here"
   feishu_base_url = "https://open.feishu.cn"
   
   # 群聊配置
   default_chat_id = "oc_default_chat_id"
   project_chat_mapping = "{\"user-service\":\"oc_user_chat\",\"payment-service\":\"oc_payment_chat\",\"default\":\"oc_default_chat\"}"
   
   # 环境变量
   lambda_environment_variables = {
     ENV = "prod"
     SERVICE_NAME = "feishu-notification-service"
   }
   ```

3. **部署服务**
   ```bash
   # 初始化Terraform
   terraform init
   
   # 查看部署计划
   terraform plan
   
   # 执行部署
   terraform apply
   ```

4. **获取接口地址**
   ```bash
   # 查看部署输出
   terraform output
   
   # 获取通知接口地址
   terraform output notification_endpoint
   
   # 获取健康检查地址
   terraform output health_check_endpoint
   ```

## 使用示例

### 1. 健康检查
```bash
curl -X GET "https://your-api-id.execute-api.cn-northwest-1.amazonaws.com.cn/prod/health"
```

### 2. 发送错误通知
```bash
curl -X POST "https://your-api-id.execute-api.cn-northwest-1.amazonaws.com.cn/prod" \
  -H "Content-Type: application/json" \
  -d '{
    "message_type": "error",
    "title": "数据库连接错误",
    "content": "无法连接到MySQL数据库，连接超时",
    "source": {
      "service_name": "user-service",
      "module_name": "database",
      "environment": "prod",
      "version": "v1.2.3"
    },
    "details": {
      "level": "ERROR",
      "timestamp": "2024-01-15T10:30:00Z",
      "error_code": "DB_TIMEOUT",
      "trace_id": "trace-123"
    },
    "target": {
      "project_name": "user-service",
      "priority": "high"
    }
  }'
```

### 3. 发送成功通知
```bash
curl -X POST "https://your-api-id.execute-api.cn-northwest-1.amazonaws.com.cn/prod" \
  -H "Content-Type: application/json" \
  -d '{
    "message_type": "success",
    "title": "部署成功",
    "content": "用户服务v1.3.0部署完成",
    "source": {
      "service_name": "user-service",
      "environment": "prod"
    },
    "target": {
      "project_name": "user-service"
    }
  }'
```

### 4. 发送到多个群聊
```bash
curl -X POST "https://your-api-id.execute-api.cn-northwest-1.amazonaws.com.cn/prod" \
  -H "Content-Type: application/json" \
  -d '{
    "message_type": "alert",
    "title": "系统告警",
    "content": "CPU使用率超过90%",
    "target": {
      "chat_ids": ["oc_chat1", "oc_chat2", "oc_chat3"]
    }
  }'
```

## 消息格式示例

发送到飞书的消息将被格式化为：

```
🚨 错误报警
标题: 数据库连接错误
服务: user-service/database
环境: prod(cn-northwest-1)
级别: ERROR
时间: 2024-01-15 18:30:00
内容: 无法连接到MySQL数据库，连接超时
追踪ID: trace-123
错误码: DB_TIMEOUT
版本: v1.2.3
```

## 项目集成

### 1. 错误监控模块集成
```go
// 发送错误通知
func sendErrorNotification(err error, context map[string]interface{}) {
    notification := map[string]interface{}{
        "message_type": "error",
        "content":      err.Error(),
        "source": map[string]interface{}{
            "service_name": "error-monitor",
            "environment":  os.Getenv("ENV"),
        },
        "details": map[string]interface{}{
            "level":     "ERROR",
            "timestamp": time.Now().Format(time.RFC3339),
            "metadata":  context,
        },
        "target": map[string]interface{}{
            "project_name": "error-monitor",
        },
    }
    
    // HTTP POST到Lambda接口
    sendToLambda(notification)
}
```

### 2. CI/CD集成
```yaml
# GitHub Actions示例
- name: Send deployment notification
  run: |
    curl -X POST "${{ secrets.FEISHU_LAMBDA_URL }}" \
      -H "Content-Type: application/json" \
      -d '{
        "message_type": "success",
        "title": "部署成功",
        "content": "应用已成功部署到生产环境",
        "source": {
          "service_name": "${{ github.repository }}",
          "environment": "prod",
          "version": "${{ github.sha }}"
        },
        "target": {
          "project_name": "ci-cd"
        }
      }'
```

### 3. 监控告警集成
```python
# Python示例
import requests
import json

def send_alert(metric_name, value, threshold):
    notification = {
        "message_type": "alert",
        "title": f"{metric_name} 告警",
        "content": f"指标 {metric_name} 当前值 {value} 超过阈值 {threshold}",
        "source": {
            "service_name": "monitoring",
            "environment": "prod"
        },
        "details": {
            "level": "WARN",
            "timestamp": datetime.now().isoformat(),
            "metadata": {
                "metric": metric_name,
                "value": value,
                "threshold": threshold
            }
        },
        "target": {
            "project_name": "monitoring",
            "priority": "high"
        }
    }
    
    response = requests.post(
        "https://your-api-id.execute-api.cn-northwest-1.amazonaws.com.cn/prod",
        headers={"Content-Type": "application/json"},
        data=json.dumps(notification)
    )
```

## 故障排除

### 1. 常见问题

**问题**: 部署失败，提示权限不足
```
解决方案: 
1. 检查AWS凭证配置
2. 确保IAM用户有Lambda、API Gateway、CloudWatch权限
3. 检查中国区域的特殊配置
```

**问题**: 飞书消息发送失败
```
解决方案:
1. 检查飞书App ID和App Secret配置
2. 确认机器人已添加到目标群聊
3. 检查群聊ID是否正确
4. 查看CloudWatch日志获取详细错误信息
```

**问题**: API Gateway调用失败
```
解决方案:
1. 检查API Gateway是否正确部署
2. 确认Lambda权限配置正确
3. 检查请求格式是否符合接口规范
```

### 2. 日志查看
```bash
# 查看Lambda函数日志
aws logs describe-log-groups --log-group-name-prefix /aws/lambda/feishu-notification-service --region cn-northwest-1

# 实时查看日志
aws logs tail /aws/lambda/feishu-notification-service --follow --region cn-northwest-1
```

### 3. 调试模式
在terraform.tfvars中设置：
```hcl
lambda_environment_variables = {
  ENV = "dev"
  LOG_LEVEL = "debug"
}
```

## 更新和维护

### 1. 更新Lambda函数
```bash
# 修改源代码后重新部署
terraform apply
```

### 2. 更新配置
```bash
# 修改terraform.tfvars后应用更改
terraform plan
terraform apply
```

### 3. 扩展功能
1. 修改`src/main.go`添加新功能
2. 更新接口文档
3. 添加相应的测试用例
4. 重新部署

## 成本优化

- **按需付费**: Lambda按实际调用次数和执行时间计费
- **免费额度**: AWS免费套餐包含每月100万次Lambda请求
- **冷启动优化**: Go语言具有较快的冷启动速度
- **内存优化**: 根据实际需求调整Lambda内存分配

## 安全考虑

1. **API认证**: 可以在API Gateway层添加API Key或JWT认证
2. **网络安全**: 使用VPC Lambda限制网络访问
3. **敏感信息**: 使用AWS Secrets Manager存储敏感配置
4. **访问控制**: 通过IAM角色最小权限原则

## 监控和告警

1. **CloudWatch指标**: 自动收集Lambda执行指标
2. **自定义指标**: 可以添加业务相关的自定义指标  
3. **告警设置**: 基于错误率、延迟等设置告警
4. **日志分析**: 使用CloudWatch Logs Insights分析日志

## 许可证

MIT License

## 贡献指南

欢迎提交Issue和Pull Request来改进这个项目。在提交代码前，请确保：

1. 代码通过所有测试
2. 遵循Go代码规范
3. 更新相关文档
4. 添加必要的测试用例
