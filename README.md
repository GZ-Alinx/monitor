# Prometheus + grafana + alertmanager 监控体系

# 环境部署
## 1. 部署Prometheus
```bash
$ wget https://github.com/prometheus/prometheus/releases/download/v3.0.1/prometheus-3.0.1.linux-amd64.tar.gz
$ tar -xf prometheus-3.0.1.linux-amd64.tar.gz -C /usr/local/
$ cd !$ && ln -sv prometheus-2.28.1.linux-amd64 prometheus
$ cd prometheus
$ mkdir rules targets





$ vim prometheus.yml
global:
  scrape_interval: 10s
  evaluation_interval: 10s
  query_log_file: 'query.log'

# 告警对接配置
alerting:
  alertmanagers:
  - static_configs:
    - targets:
      - 'localhost:9093'

# 规则引擎具体规则配置
rule_files:
  - /usr/local/prometheus/rules/*.yaml

scrape_configs:
  - job_name: 'prometheus'
    file_sd_configs:
    - files:
      - /usr/local/prometheus/targets/prometheus*.yaml

  - job_name: 'alertmanager-server'
    file_sd_configs:
    - files:
      - /usr/local/prometheus/targets/alertmanager*.yaml

  - job_name: 'nodes'
    file_sd_configs:
    - files:
      - /usr/local/prometheus/targets/nodes*.yaml



# 创建目录备用
$ mkdir /usr/local/prometheus/rules
$ mkdir /usr/local/prometheus/targets


# 服务管理文件配置
$ vim /usr/lib/systemd/system/prometheus.service
[Unit]
Description=Prometheus Node Exporter
After=network.target

[Service]
ExecStart=/usr/local/prometheus/prometheus \
    --config.file=/usr/local/prometheus/prometheus.yml \
    --web.read-timeout=5m \
    --web.max-connections=500 \
    --query.max-concurrency=2000 \
    --query.timeout=2m
[Install]
WantedBy=multi-user.target


# 开机启动
$ systemctl enable prometheus.service
$ systemctl start prometheus.service

```


## 2.部署Grafana
- 看板视图地址：  https://grafana.com/grafana/dashboards
```bash
docker run -itd -p 9800:3000 --name grafana-server-prod grafana/grafana
```





## 3.部署alertmanager
```bash
$ wget https://github.com/prometheus/alertmanager/releases/download/v0.22.2/alertmanager-0.22.2.linux-amd64.tar.gz
$ tar -xf alertmanager-0.22.2.linux-amd64.tar.gz -C /usr/local/
$ cd /usr/local/ && ln -sv alertmanager-0.22.2.linux-amd64 alertmanager
$ cd /usr/local/alertmanager && mkdir config bin data logs
$ mv alertmanager amtool bin/
$ mv alertmanager.yml conf/
$ mv alertmanager.yml config/
$ mv alertmanager amtool bin/


$ vim config/alertmanager.yml
global:
  resolve_timeout: 5m

route:
  group_by: ['alertname', 'severity']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 1h
  receiver: 'telegram'

receivers:
- name: 'telegram'
  webhook_configs:
  - url: 'http://localhost:9080/alert'



$ vim /usr/lib/systemd/system/alertmanager.service
[Unit]
Description=alertmanager
Documentation=https://prometheus.io/
After=network.target

[Service]
ExecStart=/usr/local/alertmanager/bin/alertmanager --config.file=/usr/local/alertmanager/config/alertmanager.yml --log.level=info --log.format=json --storage.path="/urs/local/alertmanager/data/"
ExecReload=/bin/kill -HUP $MAINPID
ExecStop=/bin/kill -KILL $MAINPID
Type=simple
Restart=always
TimeoutStopSec=20s

[Install]
WantedBy=multi-user.target



$ systemctl daemon-reload
$ systemctl enable alertmanager
$ systemctl start alertmanager

```


## 4.开发自己的告警服务
- 暴露端口和接口对接到alertmanager
- http://localhost:9080/alert
- *符合Prometheus的告警消息请求体格式即可转发告警信息到你想要发送的地方*
    - prometheus-alert-telegram 这个项目是对telegram开放的




## 配置资源告警规则（主机）
- 告警规则获取地址：  https://samber.github.io/awesome-prometheus-alerts
- hosts.rules.yaml
```bash
$ vim /usr/local/prometheus/rules/hosts.rules.yaml
groups:
- name: Hosts
  rules:
  - alert: HighCPUUsage
    expr: 100 - (avg by (instance) (rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 80
    for: 15s
    labels:
      severity: critical
    annotations:
      summary: "Instance {{ $labels.instance }} CPU usage is high"
      description: "CPU usage is above 80% (current value: {{ $value }})"
 
  - alert: HighMemoryUsage
    expr: node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10
    for: 15s
    labels:
      severity: warning
    annotations:
      summary: "Instance {{ $labels.instance }} memory usage is high"
      description: "Memory usage is below 10% (current value: {{ $value }})"
  
  - alert: HostOutOfMemory
    expr: (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10) * on(instance) group_left (nodename) node_uname_info{nodename=~".+"}
    for: 1m
    labels:
      severity: warning
    annotations:
      summary: Host out of memory (instance {{ $labels.instance }})
      description: "Node memory is filling up (< 10% left)\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}"

  - alert: HostMemoryUnderMemoryPressure
    expr: (rate(node_vmstat_pgmajfault[1m]) > 1000) * on(instance) group_left (nodename) node_uname_info{nodename=~".+"}
    for: 1m
    labels:
      severity: warning
    annotations:
      summary: Host memory under memory pressure (instance {{ $labels.instance }})
      description: "The node is under heavy memory pressure. High rate of major page faults\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}"

  - alert: HostOutOfDiskSpace
    expr: ((node_filesystem_avail_bytes * 100) / node_filesystem_size_bytes < 10 and ON (instance, device, mountpoint) node_filesystem_readonly == 0) * on(instance) group_left (nodename) node_uname_info{nodename=~".+"}
    for: 1m
    labels:
      severity: warning
    annotations:
      summary: Host out of disk space (instance {{ $labels.instance }})
      description: "Disk is almost full (< 10% left)\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}"

  - alert: HostHighCpuLoad
    expr: (sum by (instance) (avg by (mode, instance) (rate(node_cpu_seconds_total{mode!="idle"}[2m]))) > 0.8) * on(instance) group_left (nodename) node_uname_info{nodename=~".+"}
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Host high CPU load (instance {{ $labels.instance }})
      description: "CPU load is > 80%\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}"
 
  - alert: HostCpuHighIowait
    expr: (avg by (instance) (rate(node_cpu_seconds_total{mode="iowait"}[5m])) * 100 > 10) * on(instance) group_left (nodename) node_uname_info{nodename=~".+"}
    for: 1m
    labels:
      severity: warning
    annotations:
      summary: Host CPU high iowait (instance {{ $labels.instance }})
      description: "CPU iowait > 10%. A high iowait means that you are disk or network bound.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}"
  
  - alert: HostUnusualDiskIo
    expr: (rate(node_disk_io_time_seconds_total[1m]) > 0.5) * on(instance) group_left (nodename) node_uname_info{nodename=~".+"}
    for: 1m
    labels:
      severity: warning
    annotations:
      summary: Host unusual disk IO (instance {{ $labels.instance }})
      description: "Time spent in IO is too high on {{ $labels.instance }}. Check storage for issues.\n  VALUE = {{ $value }}\n  LABELS = {{ $labels }}"

```

## 配置监控对象
- 根据需求添加更多资源对象
- nodes-server.yaml  
- prometheus-servers.yaml
```bash
$ vim /usr/local/prometheus/targets/nodes-server.yaml  
- targets:
  - '172.31.3.43:9100'
  labels:
    environment: prod
    host: bigdata-jumpserver
- targets:
  - '172.31.12.85:9100'
  labels:
    environment: prod
    host: dp-svn-1
- targets:
  - '172.31.30.199:9100'
  labels:
    environment: prod
    host: prod-data-gateway
- targets:
  - '172.31.22.190:9100'
  labels:
    environment: prod
    host: prod-doris001
- targets:
  - '172.31.26.126:9100'
  labels:
    environment: prod
    host: prod-doris002
- targets:
  - '172.31.30.25:9100'
  labels:
    environment: prod
    host: prod-doris003
- targets:
  - '172.31.5.56:9100'
  labels:
    environment: prod
    host: prod-hadoop-cicd
- targets:
  - '172.31.4.194:9100'
  labels:
    environment: prod
    host: prod-hadoop001


$ vim /usr/local/prometheus/targets/prometheus-servers.yaml
- targets:
  - 172.31.36.241:9090
  labels:
    app: prometheus-server

```



