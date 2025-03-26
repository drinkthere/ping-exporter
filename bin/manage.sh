
#!/bin/bash
cpuidx=1
# 根据输入参数执行相应的操作
case "$1" in
    "start")
        echo "Starting Service..."
        nohup taskset -c "$cpuidx" ./ping-exporter ../config/config.json > /data/dc/ping-exporter/te.log 2>&1 &
        ;;
    "stop")
        echo "Stopping Service..."
        pkill -f "./ping-exporter ../config/config.json"
        ;;
    *)
        echo "Usage: $0 {start|stop}"
        exit 1
        ;;
esac