package storage

import (
	"context"
	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
)

type storage interface {
	Сreate(ctx context.Context, cfg *config.ServerConfig) storage
	SetMetric(context.Context, metrics.Metric) error
	SetMetricsBatch(context.Context, []metrics.Metric) error
	GetMetric(context.Context, metrics.Metric) (metrics.Metric, error)
	GetAllMetrics(context.Context) ([]metrics.Metric, error)
	GetAll(context.Context) ([]string, error)
	PingStorage(context.Context) error
	Close(context.Context)
}

func Init(ctx context.Context, conf *config.ServerConfig) {
	if conf.DatabaseAddress != "" {
		ServerStorage.storage = db.Сreate(ctx, conf)
		return
	}
	ServerStorage.storage = memS.Сreate(ctx, conf)
}

type stor struct {
	storage
}

var ServerStorage stor
