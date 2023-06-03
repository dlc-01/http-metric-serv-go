package routine

import (
	"context"
	"github.com/dlc-01/http-metric-serv-go/internal/agent/collector"
	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/go-resty/resty/v2"
	"time"
)

var (
	done   = make(chan bool)
	client = resty.New()
	cfg    *config.AgentConfig
)

func Run(cfg *config.AgentConfig) {
	//TODO тз на эти итеры было оч расплывчатое мб что то не правильно понял, тесты тоже не работают
	reportTicker := time.Duration(time.Second * time.Duration(cfg.Report))
	poolTicker := time.Duration(time.Second * time.Duration(cfg.Poll))
	chanStor := make(chan []metrics.Metric)
	go collector.CollectMetrics(context.Background(), chanStor, poolTicker)
	go sendMetrics(cfg, <-chanStor, reportTicker)

}
