# AWS权限问题解决指南

## 问题分析

您当前遇到的错误表明用户 `lambda_fullAC` 缺少以下关键权限：
1. **API Gateway权限**: `apigateway:PUT` - 创建和管理API Gateway
2. **IAM权限**: `iam:CreateRole`, `iam:CreatePolicy` - 创建Lambda执行角色和策略

## 解决方案1：添加完整权限（推荐）

### 步骤1：使用提供的策略文件
我已经为您创建了 `required-iam-policy.json` 文件，包含了所有必需的权限。

### 步骤2：让AWS管理员应用此策略
请将 `required-iam-policy.json` 文件发送给您的AWS管理员，请他们为用户 `lambda_fullAC` 创建并附加这个策略：

```bash
# AWS管理员执行以下命令：
aws iam create-policy \
    --policy-name FeishuLambdaDeployPolicy \
    --policy-document file://required-iam-policy.json \
    --region cn-northwest-1

aws iam attach-user-policy \
    --user-name lambda_fullAC \
    --policy-arn arn:aws-cn:iam::325289648576:policy/FeishuLambdaDeployPolicy \
    --region cn-northwest-1
```

### 步骤3：验证权限
权限添加后，重新尝试部署：
```bash
terraform apply
```

---

## 解决方案2：使用现有IAM角色（临时方案）

如果无法获得完整权限，可以使用现有的IAM角色。

### 修改Terraform配置使用现有角色：

```hcl
# 在 terraform.tfvars 中添加：
existing_lambda_role_arn = "arn:aws-cn:iam::325289648576:role/existing-lambda-role"
```

### 然后修改 main.tf：
```hcl
# 在 variables.tf 中添加：
variable "existing_lambda_role_arn" {
  description = "现有Lambda角色ARN"
  type        = string
  default     = ""
}

# 在 main.tf 中修改Lambda模块配置：
module "lambda_function" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "~> 4.0"

  function_name = var.lambda_function_name
  description   = "飞书通知服务Lambda函数"
  handler       = "bootstrap"
  runtime       = var.lambda_runtime
  timeout       = var.lambda_timeout

  source_path = "${path.module}/src"

  # 使用现有角色而不是创建新角色
  create_role = var.existing_lambda_role_arn == "" ? true : false
  lambda_role = var.existing_lambda_role_arn

  # 其他配置保持不变...
}
```

---

## 解决方案3：简化部署（最小权限版本）

如果权限限制很严格，我们可以创建一个无API Gateway的简化版本，只部署Lambda函数：

### 创建简化版本的main.tf：
```hcl
# 仅Lambda函数，不包含API Gateway
provider "aws" {
  region = var.aws_region
}

# 如果您有现有的Lambda执行角色，请使用它
data "aws_iam_role" "existing_lambda_role" {
  count = var.existing_lambda_role_name != "" ? 1 : 0
  name  = var.existing_lambda_role_name
}

module "lambda_function" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "~> 4.0"

  function_name = var.lambda_function_name
  description   = "飞书通知服务Lambda函数"
  handler       = "bootstrap"
  runtime       = var.lambda_runtime
  timeout       = var.lambda_timeout

  source_path = "${path.module}/src"

  # 使用现有角色
  create_role = false
  lambda_role = data.aws_iam_role.existing_lambda_role[0].arn

  environment_variables = merge(var.lambda_environment_variables, {
    FEISHU_APP_ID         = var.feishu_app_id
    FEISHU_APP_SECRET     = var.feishu_app_secret
    FEISHU_BASE_URL       = var.feishu_base_url
    DEFAULT_CHAT_ID       = var.default_chat_id
    PROJECT_CHAT_MAPPING  = var.project_chat_mapping
  })

  tags = {
    Environment = lookup(var.lambda_environment_variables, "ENV", "dev")
    Project     = "feishu-notification-service"
    Service     = "notification"
  }
}
```

### 在variables.tf中添加：
```hcl
variable "existing_lambda_role_name" {
  description = "现有Lambda角色名称"
  type        = string
  default     = ""
}
```

---

## 解决方案4：使用AWS CLI直接部署Lambda

如果Terraform权限受限，可以直接使用AWS CLI部署：

### 步骤1：打包Lambda函数
```bash
cd src
chmod +x build.sh
./build.sh
zip ../lambda-function.zip bootstrap
cd ..
```

### 步骤2：创建Lambda函数（如果有Lambda权限）
```bash
aws lambda create-function \
    --function-name feishu-notification-service \
    --runtime provided.al2 \
    --role arn:aws-cn:iam::325289648576:role/your-existing-lambda-role \
    --handler bootstrap \
    --zip-file fileb://lambda-function.zip \
    --timeout 60 \
    --environment Variables='{
        FEISHU_APP_ID=cli_a8028d701ebad00c,
        FEISHU_APP_SECRET=PMHEgtwRpyQicbsnMbuCyhdeIpuat2cX,
        FEISHU_BASE_URL=https://open.feishu.cn,
        PROJECT_CHAT_MAPPING="{\"default\":\"\"}"
    }' \
    --region cn-northwest-1
```

### 步骤3：直接调用Lambda函数测试
```bash
# 创建测试数据
cat > test-event.json << 'EOF'
{
  "httpMethod": "POST",
  "path": "/",
  "body": "{\"message_type\":\"info\",\"content\":\"测试消息\",\"target\":{\"chat_ids\":[\"your_chat_id\"]}}"
}
EOF

# 调用Lambda函数
aws lambda invoke \
    --function-name feishu-notification-service \
    --payload file://test-event.json \
    --response.json \
    --region cn-northwest-1
```

---

## 推荐的解决顺序

1. **首选**: 联系AWS管理员添加完整权限（解决方案1）
2. **次选**: 使用现有IAM角色（解决方案2）
3. **备选**: 简化部署，去掉API Gateway（解决方案3）
4. **最后**: 直接使用AWS CLI部署（解决方案4）

## 如何确认权限是否足够

在尝试部署前，可以测试权限：

```bash
# 测试Lambda权限
aws lambda list-functions --region cn-northwest-1

# 测试API Gateway权限
aws apigateway get-rest-apis --region cn-northwest-1

# 测试IAM权限
aws iam list-roles --region cn-northwest-1
```

如果这些命令都能成功执行，那么权限应该足够了。

## 联系管理员时的说明

如果需要联系AWS管理员，可以这样说明：

> 您好，我需要部署一个飞书通知Lambda函数，需要为用户 lambda_fullAC 添加以下权限：
> 1. Lambda函数的完整权限
> 2. API Gateway的创建和管理权限  
> 3. IAM角色和策略的创建权限
> 4. CloudWatch日志权限
> 
> 我已经准备了完整的权限策略文件 (required-iam-policy.json)，请帮忙创建并附加到我的用户。

这样应该能解决您遇到的权限问题。建议优先尝试解决方案1，如果不行再考虑其他方案。
