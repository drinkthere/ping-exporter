package config

import (
	"encoding/json"
	"go.uber.org/zap/zapcore"
	"os"
)

type Config struct {
	SourceIPs []string
	TargetIPs []string

	// 日志配置
	LogLevel zapcore.Level
	LogPath  string

	// HTTP服务的端口
	Port int
	// 上报Prometheus时，exporter的Tag
	PrometheusTag string
	// 每组IP ping的次数
	PingTimes int
}

func LoadConfig(filename string) *Config {
	config := new(Config)
	reader, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	// 加载配置
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}
	return config
}
