package main

import (
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/agent/flags"
	"github.com/dlc-01/http-metric-serv-go/internal/agent/metrics"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
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
			//urls := m.GenerateURLMetrics(flags.ServerAddress)
			//for _, url := range urls {
			//	client.R().SetHeader("Content-Type", "text/plain").Post(url)
			//}

			for metric, value := range m.Gauge {
				request := storage.Metrics{
					ID:    metric,
					MType: "gauge",
					Value: &value,
				}
				client.R().SetHeader("Content-Type", "application/json").
					SetBody(request).
					Post(fmt.Sprintf("http://%s/update/", flags.ServerAddress))
			}
			for metric, value := range m.Counter {
				request := storage.Metrics{
					ID:    metric,
					MType: "counter",
					Delta: &value,
				}
				client.R().SetHeader("Content-Type", "application/json").
					SetBody(request).
					Post(fmt.Sprintf("http://%s/update/", flags.ServerAddress))
			}

			/*
				Как лучше сделать?
				Как ниже коммит или как используется?
			*/

			//requests := m.GenerateStructMetrics()
			//for _, request := range requests {
			//	client.R().SetHeader("Content-Type", "application/json").
			//		SetBody(request).
			//		Post(fmt.Sprintf("http://%s/update/", flags.ServerAddress))
			//}

		case <-t2.C:
			m.Check()
		case <-term:
			fmt.Printf("%+v\n", "termination")
			running = false
		}
	}

}
