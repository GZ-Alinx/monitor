#!/bin/bash

# 检查Docker是否已安装
if ! command -v docker &> /dev/null
then
    echo "Docker未安装，请先安装Docker"
    yum install -y docker
    systemctl enable docker
    systemctl start docker
fi

# 拉取node exporter镜像
echo "正在拉取Prometheus node exporter镜像..."
docker pull prom/node-exporter:latest

# 运行node exporter容器
echo "正在启动node exporter容器..."
docker run -d \
  --name="node-exporter" \
  --restart=always \
  --net="host" \
  --pid="host" \
  -v "/:/host:ro,rslave" \
  prom/node-exporter:latest \
  --path.rootfs=/host

# 检查容器状态
echo "Node exporter已启动，状态如下："
docker ps -f name=node-exporter

echo "安装完成！Node exporter正在运行，可以通过http://localhost:9100/metrics访问"
