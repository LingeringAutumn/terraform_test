#!/bin/bash

# 构建Go Lambda函数的脚本

echo "开始构建Go Lambda函数..."

# 进入源代码目录
cd "$(dirname "$0")"

# 设置Go环境变量
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

# 下载依赖
echo "下载Go模块依赖..."
go mod tidy

# 编译Go代码
echo "编译Go代码..."
go build -o main main.go

if [ $? -eq 0 ]; then
    echo "Go Lambda函数构建成功！"
    echo "生成的可执行文件: main"
    ls -la main
else
    echo "构建失败！"
    exit 1
fi
