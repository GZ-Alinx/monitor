#!/bin/bash

# 设置程序名称和路径
APP_NAME="prometheus-alert-telegram-linux-1.0.0"
APP_PATH="/usr/local/prometheus-alert-telegram/"
LOG_PATH="/usr/local/prometheus-alert-telegram/$APP_NAME.log"
export GIN_MODE=release

# 启动程序
function start_app {
    if pgrep -f "$APP_NAME" > /dev/null; then
        echo "$APP_NAME is already running."
        return 1
    else
        echo "Starting $APP_NAME..."
              "$APP_PATH/$APP_NAME" > "$LOG_PATH" 2>&1 &
        if [ $? -eq 0 ]; then
            echo "$APP_NAME started with PID $!"
            return 0 
        else
            echo "Failed to start $APP_NAME."
            return 1
        fi
    fi
}

# 停止程序
function stop_app {
    if pgrep -f "$APP_NAME" > /dev/null; then
        echo "Stopping $APP_NAME..."
        pkill -f "$APP_NAME"
        if [ $? -eq 0 ]; then
            echo "$APP_NAME stopped."
            return 0
        else
            echo "Failed to stop $APP_NAME."
            return 1
        fi
    else
        echo "$APP_NAME is not running."
        return 1
    fi
}

# 重启程序
function restart_app {
    stop_app
    sleep 1
    start_app
}

# 打印使用帮助
function usage {
    echo "Usage: $0 {start|stop|restart}"
    exit 1
}

# 根据输入的参数执行相应的操作
case "$1" in
    start)
        start_app
        ;;
    stop)
        stop_app
        ;;
    restart)
        restart_app
        ;;
    *)
        usage
        ;;
esac