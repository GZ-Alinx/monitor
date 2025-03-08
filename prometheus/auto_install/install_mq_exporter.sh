#!/bin/bash

# 检查当前用户是否为root
if [ "$EUID" -ne 0 ]; then
  echo "请以root权限运行此脚本"
  exit 1
fi

# 检查Docker是否已安装
if ! command -v docker &> /dev/null
then
  echo "Docker 未安装，正在安装..."
  yum update -y
  yum install -y yum-utils
  yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
  yum install -y docker-ce docker-ce-cli containerd.io
  systemctl start docker
  systemctl enable docker
else
  echo "Docker 已安装"
fi

# 确认Docker启动
if ! systemctl is-active --quiet docker
then
  echo "启动Docker..."
  systemctl start docker
  systemctl enable docker
fi

# 拉取 RabbitMQ Exporter 镜像
echo "正在拉取 RabbitMQ Exporter Docker 镜像..."
docker pull kbudde/rabbitmq-exporter

# 启动 RabbitMQ Exporter 使用Host网络模式
echo "启动 RabbitMQ Exporter，使用Host网络模式..."
docker run -d \
  --name rabbitmq-exporter \
  --net=host \
  -e RABBIT_URL="http://localhost:15672" \
  -e RABBIT_USER="guest" \
  -e RABBIT_PASSWORD="guest" \
  kbudde/rabbitmq-exporter

# 校验是否安装成功
if docker ps | grep -q rabbitmq-exporter; then
  echo "RabbitMQ Exporter 安装并运行成功"
else
  echo "RabbitMQ Exporter 安装失败，请检查日志"
  exit 1
fi