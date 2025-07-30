# 使用Terraform部署AWS Lambda

本项目演示如何使用Terraform将AWS Lambda函数部署到AWS。

## 前提条件

在开始之前，请确保您已完成以下步骤：

1. **安装Terraform**
   - 前往[Terraform官网](https://www.terraform.io/downloads.html)下载适合您操作系统的版本。
   - 解压下载的文件并将其路径添加到系统的环境变量中。
   - 验证安装是否成功：
     ```bash
     terraform -v
     ```

2. **安装AWS CLI**
   - 前往[AWS CLI官网](https://aws.amazon.com/cli/)下载并安装AWS CLI。
   - 验证安装是否成功：
     ```bash
     aws --version
     ```

3. **配置AWS凭证**
   - 运行以下命令并根据提示输入您的AWS访问密钥和区域：
     ```bash
     aws configure
     ```
   - 验证凭证是否正确配置：
     ```bash
     aws sts get-caller-identity --region cn-northwest-1 --endpoint-url https://sts.cn-northwest-1.amazonaws.com.cn
     ```

## 项目结构

```
aws-lambda/
├── README.md             # 项目文档
├── main.tf               # 主Terraform配置文件，使用AWS Lambda模块
├── variables.tf          # Terraform变量定义
├── terraform.tfvars      # Terraform变量值配置文件（自定义配置在这里）
├── terraform.tfvars.example # 变量配置示例文件
├── outputs.tf            # Terraform输出定义
├── versions.tf           # Terraform版本要求配置
└── src/                  # Lambda函数源代码
    └── index.js          # 示例Node.js Lambda函数
```

## 修改Lambda函数代码

1. 打开`src/index.js`文件。
2. 替换为您自己的Lambda函数代码。例如：
   ```javascript
   exports.handler = async (event) => {
       return {
           statusCode: 200,
           body: JSON.stringify({ message: "Hello, Custom Lambda!" }),
       };
   };
   ```

## 自定义部署配置

通过修改`terraform.tfvars`文件（不是variables.tf）来自定义您的部署：

1. 如果`terraform.tfvars`文件不存在，可以复制示例文件：
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   ```

2. 编辑`terraform.tfvars`文件，根据您的需求修改变量值：
   ```
   # 部署AWS Lambda函数的区域
   aws_region = "cn-northwest-1"
   
   # Lambda函数的名称
   lambda_function_name = "my-custom-function"
   
   # Lambda函数的运行时环境（nodejs18.x, python3.9, java11等）
   lambda_runtime = "python3.9"
   
   # Lambda函数的超时时间（秒）
   lambda_timeout = 60
   
   # Lambda函数的环境变量
   lambda_environment_variables = {
     ENV = "prod"
     API_KEY = "your-api-key"
     DEBUG = "false"
   }
   ```

## 部署步骤

1. **初始化Terraform**
   - 在项目根目录运行以下命令：
     ```bash
     terraform init
     ```

2. **查看部署计划**
   - 运行以下命令以查看Terraform将创建的资源：
     ```bash
     terraform plan
     ```

3. **部署资源**
   - 运行以下命令以实际部署资源：
     ```bash
     terraform apply
     ```
   - 系统会要求确认，输入`yes`继续。

4. **验证Lambda函数**
   - 部署完成后，您可以使用以下命令调用Lambda函数：
     ```bash
     aws lambda invoke --function-name $(terraform output -raw lambda_function_name) --payload '{}' response.json --region $(terraform output -raw aws_region) --endpoint-url https://lambda.$(terraform output -raw aws_region).amazonaws.com.cn
     ```
   - 或者使用terraform输出的命令示例：
     ```bash
     $(terraform output -raw lambda_invoke_command)
     ```
   - 检查`response.json`文件以查看函数的响应。

## 清理资源

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

## 变量说明

项目中的主要配置变量（在terraform.tfvars中设置）：

- `aws_region`: AWS区域（默认为cn-northwest-1）
- `lambda_function_name`: Lambda函数名称
- `lambda_runtime`: Lambda运行时环境
- `lambda_timeout`: Lambda函数超时时间（秒）
- `lambda_environment_variables`: Lambda函数环境变量

## 故障排除

如遇问题，请检查：

1. AWS凭证是否正确配置。
2. Terraform版本是否兼容(需要1.0.0以上)。
3. Lambda函数代码是否有语法错误。
4. 是否正确初始化了模块依赖。
5. 确保terraform.tfvars文件中的变量设置正确。