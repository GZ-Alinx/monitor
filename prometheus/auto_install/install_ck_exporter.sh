#!/bin/bash

set -e

# 配置参数
CLICKHOUSE_EXPORTER_VERSION="latest"
CLICKHOUSE_HOST="127.0.0.1"
CLICKHOUSE_PORT="9000"
CLICKHOUSE_USER="default"
CLICKHOUSE_PASSWORD=""

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo "Docker 未安装，正在安装..."
    sudo yum install -y yum-utils
    sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
    sudo yum install -y docker-ce docker-ce-cli containerd.io
    echo "Docker 安装完成"
else
    echo "Docker 已安装"
fi

# 设置 Docker 开机启动并启动服务
sudo systemctl enable docker
sudo systemctl start docker

# 运行 ClickHouse-Exporter 容器
echo "启动 ClickHouse-Exporter 容器..."
docker run -d --restart=always --network=host \
    --name clickhouse-exporter \
    -e CLICKHOUSE_EXPORTER_LISTEN_PORT=9116 \
    -e CLICKHOUSE_HOST=$CLICKHOUSE_HOST \
    -e CLICKHOUSE_PORT=$CLICKHOUSE_PORT \
    -e CLICKHOUSE_USER=$CLICKHOUSE_USER \
    -e CLICKHOUSE_PASSWORD=$CLICKHOUSE_PASSWORD \
    docker.io/f1yegor/clickhouse-exporter:$CLICKHOUSE_EXPORTER_VERSION

echo "ClickHouse-Exporter 安装完成并运行中，监听端口 9116"