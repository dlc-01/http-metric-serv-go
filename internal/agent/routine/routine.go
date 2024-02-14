package routine

import (
	"context"
	"github.com/dlc-01/http-metric-serv-go/internal/agent/grpcclient"
	_ "net/http/pprof"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/dlc-01/http-metric-serv-go/internal/agent/collector"
	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
)

const metricsChanSize = 1000

var (
	client  = resty.New()
	metrisC = make(chan []metrics.Metric, metricsChanSize)
)

// Run  — function use for starting goroutine that collect metric and that it to the server.
func Run(ctx context.Context, cfg *config.AgentConfig) {

	poolTicker := time.NewTicker(time.Second * time.Duration(cfg.Poll))
	go collector.CollectMetricsRuntime(ctx, metrisC, poolTicker)
	go collector.CollectMetricsGopsutil(ctx, metrisC, poolTicker)
	if cfg.GRPC {
		go grpcclient.SendMetricsViaGrpc(cfg, metrisC)
	} else {
		go sendMetrics(cfg, metrisC)
	}

}

func Shutdown() {
	close(metrisC)
}
