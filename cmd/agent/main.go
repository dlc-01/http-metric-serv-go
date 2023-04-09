package main

import (
	"github.com/dlc-01/http-metric-serv-go/internal/agent/metrics"
	"github.com/go-resty/resty/v2"
	"time"
)

func main() {
	client := resty.New()

	metrics := new(metrics.MemMetrics)
	metrics.Init()
	timeCounter := 1

	for {
		metrics.Check()
		time.Sleep(2 * time.Second)
		timeCounter++
		if timeCounter%5 == 0 {
			urls := metrics.GenerateURLMetrics()
			for _, url := range urls {
				client.R().SetHeader("Content-Type", "text/plain").Post(url)
			}
		}

	}
}
