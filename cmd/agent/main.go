package main

import (
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/agent/flagsOs"
	"github.com/dlc-01/http-metric-serv-go/internal/agent/metrics"
	"github.com/go-resty/resty/v2"

	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	flagsOs.ParseFlagsOs()

	client := resty.New()

	m := metrics.Init()

	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGINT, syscall.SIGTERM)
	t1 := time.NewTicker(time.Second * time.Duration(flagsOs.Report))
	t2 := time.NewTicker(time.Second * time.Duration(flagsOs.Poll))
	running := true

	for running {
		select {
		case <-t1.C:
			//urls := m.GenerateURLMetrics(flagsOs.ServerAddress)
			//for _, url := range urls {
			//	client.R().SetHeader("Content-Type", "text/plain").Post(url)
			//}
			requests := m.GenerateStructMetrics()
			for _, request := range requests {
				client.R().SetHeader("Content-Type", "application/json").
					SetBody(request).
					Post(fmt.Sprintf("http://%s/update/", flagsOs.ServerAddress))

			}
		case <-t2.C:
			m.Check()
		case <-term:
			fmt.Printf("%+v\n", "termination")
			running = false
		}
	}

}
