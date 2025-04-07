# 使用Terraform部署AWS Lambda

本项目演示如何使用Terraform将Lambda函数部署到AWS。

## 前提条件

- [AWS账户](https://aws.amazon.com/)
- [AWS CLI](https://aws.amazon.com/cli/)已配置
- [Terraform](https://www.terraform.io/downloads.html)已安装(版本 >= 1.0.0)

## 项目结构

```
aws-lambda/
├── README.md             # 项目文档
├── main.tf               # 主Terraform配置文件
├── variables.tf          # Terraform变量定义
├── outputs.tf            # Terraform输出定义
└── src/                  # Lambda函数源代码
    └── index.js          # 示例Node.js Lambda函数
```

## 配置步骤

1. **配置AWS凭证**

   确保您已经配置好AWS凭证。您可以通过运行以下命令配置：

   ```bash
   aws configure
   ```

2. **初始化Terraform**

   ```bash
   terraform init
   ```

3. **查看部署计划**

   ```bash
   terraform plan
   ```

4. **部署资源**

   ```bash
   terraform apply
   ```

   系统会要求确认，输入`yes`继续。

5. **清理资源**

   当您不再需要这些资源时，可以通过以下命令删除它们：

   ```bash
   terraform destroy
   ```

## AWS中国区域特别说明

如果您使用的是AWS中国区域(cn-north-1或cn-northwest-1)，请注意：

1. 确保使用的是中国区域专用账户和凭证
2. 服务终端节点和ARN格式有所不同，已在配置中自动调整
3. 如遇凭证问题，请使用以下命令验证凭证是否正确配置：

```bash
aws sts get-caller-identity --region cn-northwest-1 --endpoint-url https://sts.cn-northwest-1.amazonaws.com.cn
```

## Lambda函数说明

本项目部署了一个简单的Node.js Lambda函数，当触发时返回一条"Hello, World!"消息。您可以根据自己的需求修改`src/index.js`文件。

## 资源说明

本Terraform配置创建以下AWS资源：

- Lambda函数
- IAM角色和策略（允许Lambda函数执行和记录日志）
- CloudWatch日志组（用于存储Lambda函数日志）

## 自定义配置

您可以通过修改`variables.tf`文件中的变量来自定义部署：

- `aws_region`: AWS区域
- `lambda_function_name`: Lambda函数名称
- `lambda_runtime`: Lambda运行时环境

## 故障排除

如遇问题，请检查：

1. AWS凭证是否正确配置
2. Terraform版本是否兼容
3. Lambda函数代码是否有语法错误