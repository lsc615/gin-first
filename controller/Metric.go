package controller

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func TestMetric() {
	var (
		// 自定义一个HTTP请求持续时间的Histogram
		success = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "success",
				Help: "A histogram of the HTTP request durations in seconds.",
			},
		)
	)

	// Prometheus的默认注册表
	registry := prometheus.NewRegistry()
	// 在init函数中注册自定义的metrics
	registry.MustRegister(success)

	success.Set(1)

	// gin 等 web 框架中可以参考对应框架加载示例
	// http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{EnableOpenMetrics: true}))

	// 使用单独的goroutine来启动Prometheus metrics的HTTP服务器
	go func() {
		http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
		if err := http.ListenAndServe(":9090", nil); err != nil {
			panic(err)
		}
	}()
}
