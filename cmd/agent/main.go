package main

import (
	"flag"
	"github.com/dlc-01/http-metric-serv-go/internal/agent/metrics"
	"github.com/go-resty/resty/v2"
	"os"
	"strconv"
	"time"
)

var (
	serverAddress string
	report        int
	poll          int
)

func parseFlagsOs() {
	flag.StringVar(&serverAddress, "a", "localhost:8080", "server address")
	flag.IntVar(&report, "r", 10, "report interval")
	flag.IntVar(&poll, "p", 2, "poll interval")
	flag.Parse()
	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		serverAddress = envServerAddress
	}
	if envReport := os.Getenv("REPORT_INTERVAL"); envReport != "" {
		intReport, err := strconv.ParseInt(envReport, 10, 32)
		if err != nil {
			panic(err)
		}
		report = int(intReport)
	}
	if envPoll := os.Getenv("POLL_INTERVAL"); envPoll != "" {
		intPoll, err := strconv.ParseInt(envPoll, 10, 32)
		if err != nil {
			panic(err)
		}
		poll = int(intPoll)
	}

}

func main() {
	parseFlagsOs()
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
