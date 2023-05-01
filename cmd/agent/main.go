package main

import (
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/agent/flags"
	"github.com/dlc-01/http-metric-serv-go/internal/agent/metrics"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers/url"
	"github.com/go-resty/resty/v2"

	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	flags.ParseFlagsOs()

	client := resty.New()

	m := metrics.Init()

	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGINT, syscall.SIGTERM)
	t1 := time.NewTicker(time.Second * time.Duration(flags.Report))
	t2 := time.NewTicker(time.Second * time.Duration(flags.Poll))
	running := true

	for running {
		select {
		case <-t1.C:

			for metric, value := range m.Gauge {
				request := m.GenerateRequestBody(url.GaugeTypeName, metric, 0, value)
				client.R().SetHeader("Content-Encoding", "gzip").
					SetHeader("Accept-Encoding", "gzip").
					SetHeader("Content-Type", "application/json").
					SetBody(request).
					Post(fmt.Sprintf("http://%s/update/", flags.ServerAddress))
			}
			for metric, value := range m.Counter {
				request := m.GenerateRequestBody(url.CounterTypeName, metric, value, 0)
				client.R().SetHeader("Content-Encoding", "gzip").
					SetHeader("Accept-Encoding", "gzip").
					SetHeader("Content-Type", "application/json").
					SetBody(request).
					Post(fmt.Sprintf("http://%s/update/", flags.ServerAddress))
			}
		case <-t2.C:
			m.Check()
		case <-term:
			fmt.Printf("%+v\n", "termination")
			running = false
		}
	}

}
