#!/bin/bash

APP_NAME="prometheus-alert-telegram"
VERSION="1.0.0"
OUTPUT_DIR="./"
TEST_DIR="./test"

# 确保构建目录存在
mkdir -p $OUTPUT_DIR

# 运行测试
echo "Running tests..."
go test -v $TEST_DIR
if [ $? -ne 0 ]; then
    echo "Tests failed, aborting build"
    exit 1
fi

# 构建 Linux 平台
echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o $OUTPUT_DIR/$APP_NAME-linux-$VERSION main.go

# 构建 Windows 平台
echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -o $OUTPUT_DIR/$APP_NAME-windows-$VERSION.exe main.go

echo "Build completed. Artifacts are in the $OUTPUT_DIR directory."
