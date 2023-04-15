package main

import (
	"flag"
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/agent/metrics"
	"github.com/go-resty/resty/v2"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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
			log.Fatalf("cannot parse REPORT_INTERVAL: %v", err)
		}
		report = int(intReport)
	}

	if envPoll := os.Getenv("POLL_INTERVAL"); envPoll != "" {
		intPoll, err := strconv.ParseInt(envPoll, 10, 32)
		if err != nil {
			log.Fatalf("cannot parse POLL_INTERVAL: %v", err)
		}
		poll = int(intPoll)
	}
}

func main() {
	parseFlagsOs()

	client := resty.New()

	m := metrics.Init()

	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGINT, syscall.SIGTERM)
	t1 := time.NewTicker(time.Second * time.Duration(report))
	t2 := time.NewTicker(time.Second * time.Duration(poll))
	running := true

	for running {
		select {
		case <-t1.C:
			urls := m.GenerateURLMetrics(serverAddress)
			for _, url := range urls {
				client.R().SetHeader("Content-Type", "text/plain").Post(url)
			}
		case <-t2.C:
			m.Check()
		case <-term:
			fmt.Printf("%+v\n", "termination")
			running = false
		}
	}

}
