package main

import (
	"flag"
	"github.com/dlc-01/http-metric-serv-go/internal/agent/metrics"
	"github.com/go-resty/resty/v2"
	"time"
)

var (
	serverAddress string
	report        int
	poll          int
)

func parseFlags() {
	flag.StringVar(&serverAddress, "a", "localhost:8080", "server address")
	flag.IntVar(&report, "r", 10, "report interval")
	flag.IntVar(&poll, "p", 2, "poll interval")
	flag.Parse()
}

func main() {
	parseFlags()
	pollInterval := time.Duration(poll) * time.Second
	reportInterval := time.Duration(report) * time.Second
	client := resty.New()

	metrics := new(metrics.MemMetrics)
	metrics.Init()
	timeCounter := 0

	for {
		metrics.Check()
		time.Sleep(pollInterval)
		timeCounter++
		if timeCounter == int(reportInterval/pollInterval) {
			urls := metrics.GenerateURLMetrics(serverAddress)
			for _, url := range urls {
				client.R().SetHeader("Content-Type", "text/plain").Post(url)
			}
		}

	}
}
