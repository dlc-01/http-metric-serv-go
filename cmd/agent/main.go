package main

import (
	"github.com/dlc-01/http-metric-serv-go/internal/agent/metrics"
	"log"
	"net/http"
	"time"
)

func main() {
	host := "http://localhost:8080"
	metrics := new(metrics.MemMetrics)
	metrics.Init()
	timeCounter := 1

	for {
		metrics.Check()
		time.Sleep(2 * time.Second)
		timeCounter++
		if timeCounter%5 == 0 {
			urls := metrics.GenerateUrlMetrics(host)
			for _, url := range urls {
				req, err := http.NewRequest(http.MethodPost, url, nil)
				if err != nil {
					log.Println(err)
				}
				client := &http.Client{}
				req.Header.Set("Content-Type", "text/plain")
				resp, err := client.Do(req)
				if err != nil {
					log.Println(err)
				}
				resp.Body.Close()
			}
		}

	}
}
