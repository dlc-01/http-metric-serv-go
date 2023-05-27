package storage

import (
	"context"
	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
)

type Storage interface {
	Сreate(ctx context.Context, cfg *config.ServerConfig) Storage
	SetMetric(context.Context, metrics.Metric) error
	SetMetricsBatch(context.Context, []metrics.Metric) error
	GetMetric(context.Context, metrics.Metric) (metrics.Metric, error)
	GetAllMetrics(context.Context) ([]metrics.Metric, error)
	GetAll(context.Context) ([]string, error)
	PingStorage(context.Context) error
	Close(context.Context)
}

func Init(ctx context.Context, conf *config.ServerConfig) Storage {
	if conf.DatabaseAddress != "" {

		return db.Сreate(ctx, conf)
	}

	return memS.Сreate(ctx, conf)
}

type Stor struct {
	Storage
}

var ServerStorage Stor
