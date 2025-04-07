# 使用Terraform模块部署AWS Lambda

本项目演示如何使用Terraform的AWS Lambda模块将Lambda函数部署到AWS。

## 前提条件

- [AWS账户](https://aws.amazon.com/)
- [AWS CLI](https://aws.amazon.com/cli/)已配置
- [Terraform](https://www.terraform.io/downloads.html)已安装(版本 >= 1.0.0)

## 项目结构

```
aws-lambda/
├── README.md             # 项目文档
├── main.tf               # 主Terraform配置文件，使用AWS Lambda模块
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
2. 配置已针对中国区域的特殊服务终端节点进行了优化
3. 如遇凭证问题，请使用以下命令验证凭证是否正确配置：

```bash
aws sts get-caller-identity --region cn-northwest-1 --endpoint-url https://sts.cn-northwest-1.amazonaws.com.cn
```

## 使用的Terraform模块

本项目使用了[terraform-aws-modules/lambda/aws](https://registry.terraform.io/modules/terraform-aws-modules/lambda/aws/latest)模块，该模块提供了配置AWS Lambda函数的全面解决方案，包括以下功能：

- 创建Lambda函数
- 配置IAM角色和策略
- 设置CloudWatch日志
- 管理函数代码打包和部署
- 支持环境变量配置

## 自定义配置

您可以通过修改`variables.tf`文件中的变量来自定义部署：

- `aws_region`: AWS区域
- `lambda_function_name`: Lambda函数名称
- `lambda_runtime`: Lambda运行时环境

## 故障排除

如遇问题，请检查：

1. AWS凭证是否正确配置
2. Terraform版本是否兼容(需要1.0.0以上)
3. Lambda函数代码是否有语法错误
4. 是否正确初始化了模块依赖