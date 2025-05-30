#!/bin/bash

dst=/data/moods/bin
filename=./delicate-moods-api
port=10086

cd $dst

if [ -e $filename ]; then

  chmod +x $filename && nohup $filename -conf ../config/api_prod.yml &
  if [ $? -eq 0 ]; then
    sleep 1
    lsof -i:$port | grep -v PID | awk '{print $2}'
    echo "启动成功！"
  else
    echo "启动失败！"
  fi

else
  echo "文件不存在！"
fi
