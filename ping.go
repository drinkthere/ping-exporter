package main

import (
	"fmt"
	"github.com/drinkthere/ping"
	"ping-exporter/config"
	"ping-exporter/utils/logger"
	"sync"
)

func LoopPingAws(cfg *config.Config) string {

	var fastestLatency = float64(-1)
	var fastestSource, fastestTarget string

	// Use a wait group for concurrency
	var wg sync.WaitGroup
	mutex := &sync.Mutex{} // For protecting shared variables

	for _, sourceIP := range cfg.SourceIPs {
		for _, targetIP := range cfg.TargetIPs {
			wg.Add(1)

			go func(src, tgt string) {
				defer wg.Done()

				// Ping the source-target combination
				latency := pingWithMinimum(src, tgt, cfg.PingTimes)
				logger.Info("Minimum latency from %s to %s: %.5fms", src, tgt, latency)
				// Update the fastest latency if applicable
				mutex.Lock()
				defer mutex.Unlock()
				if latency != -1 && (fastestLatency == -1 || latency < fastestLatency) {
					fastestLatency = latency
					fastestSource = src
					fastestTarget = tgt
				}
			}(sourceIP, targetIP)
		}
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Return Prometheus-formatted metrics
	if fastestLatency != -1 {
		return fmt.Sprintf("%s{source=\"%s\", target=\"%s\"} %.4f\n",
			cfg.PrometheusTag, fastestSource, fastestTarget, fastestLatency)
	}
	return fmt.Sprintf("%s{source=\"none\", target=\"none\"} -1\n", cfg.PrometheusTag)
}

// pingWithMinimum pings a source and target multiple times and returns the minimum latency.
func pingWithMinimum(sourceIP, targetIP string, count int) float64 {
	minLatency := float64(-1)

	// 创建 pinger 实例
	pinger, err := ping.NewPinger(targetIP)
	if err != nil {
		logger.Error("Bind Target IP %s Failed %+v", targetIP, err)
		return minLatency
	}
	pinger.SetSource(sourceIP)
	pinger.Count = count // 设置 Ping 次数为 10
	err = pinger.Run()   // 阻塞直到完成
	if err != nil {
		logger.Error("Ping Source IP %s -> Target IP %s Failed %+v", sourceIP, targetIP, err)
		return minLatency
	}
	stats := pinger.Statistics() // 获取 Ping 统计信息
	return float64(stats.MinRtt.Microseconds()) / 1000
}
