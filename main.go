package main

import (
	"fmt"
	"net/http"
	"os"
	"ping-exporter/config"
	"ping-exporter/utils/logger"
	"runtime"
)

var globalConfig config.Config

func main() {
	runtime.GOMAXPROCS(1)
	// 参数判断
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s config_file\n", os.Args[0])
		os.Exit(1)
	}

	// 加载配置文件
	globalConfig = *config.LoadConfig(os.Args[1])

	// 设置日志级别, 并初始化日志
	logger.InitLogger(globalConfig.LogPath, globalConfig.LogLevel)

	// Define the /metrics endpoint
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics := LoopPingAws(&globalConfig)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(metrics))
	})

	// Start the HTTP server
	fmt.Printf("Ping Exporter running on http://localhost:%d/metrics\n", globalConfig.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", globalConfig.Port), nil)
}
