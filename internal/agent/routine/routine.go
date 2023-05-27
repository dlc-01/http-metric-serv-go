package routine

import (
	"context"
	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/go-resty/resty/v2"
	"time"
)

var (
	done   = make(chan bool)
	client = resty.New()
	cfg    *config.AgentConfig
)

type stor struct {
	storage.Storage
}

var agent stor

func Run(cfg *config.AgentConfig, s storage.Storage) {
	agent.Storage = s
	reportTicker := time.NewTicker(time.Second * time.Duration(cfg.Report))
	poolTicker := time.NewTicker(time.Second * time.Duration(cfg.Poll))
	running := true

	for running {
		select {
		case <-reportTicker.C:
			logging.Info("report")
			if err := agent.sendMetrics(cfg.ServerAddress); err != nil {
				logging.Errorf("cannot send metrics: %s", err)
			}
		case <-poolTicker.C:
			logging.Info("collect")
			agent.CollectMetrics(context.TODO())
		case <-done:
			running = false
		}
	}
}

func Shutdown() {
	done <- true
}
