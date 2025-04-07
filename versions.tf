terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"  # 使用与AWS中国区域兼容的稳定版本
    }
  }
  required_version = ">= 1.0.0"
}
