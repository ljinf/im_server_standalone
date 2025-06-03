#!/bin/bash

PORT=10086

PID=$(lsof -i:$PORT | grep -v PID | awk '{print $2}')

# 检查是否找到了PID
if [ -z "$PID" ]; then
    echo "没有找到使用端口$PORT的进程。"
else
    # 使用kill命令发送SIGTERM信号给进程，请求其终止
    # 如果进程不响应SIGTERM，你可以使用kill -9 $PID来强制终止它
    echo "正在终止PID为$PID的进程..."
    kill $PID

    # 等待一小段时间以确认进程是否已终止
    sleep 1

    # 检查进程是否仍然在运行
    if kill -0 $PID 2>/dev/null; then
        echo "进程$PID没有响应终止请求，尝试使用SIGKILL..."
        kill -9 $PID
    else
        echo "进程$PID已成功终止。"
    fi
fi
