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

func Init(ctx context.Context, conf *config.ServerConfig) {
	if conf.DatabaseAddress != "" {

		serverStorage = db.Сreate(ctx, conf)
		return
	}

	serverStorage = memS.Сreate(ctx, conf)
}

var serverStorage Storage

func SetMetric(ctx context.Context, metric metrics.Metric) error {
	return serverStorage.SetMetric(ctx, metric)
}

func SetMetricsBatch(ctx context.Context, metric []metrics.Metric) error {
	return serverStorage.SetMetricsBatch(ctx, metric)
}

func GetMetric(ctx context.Context, metric metrics.Metric) (metrics.Metric, error) {
	return serverStorage.GetMetric(ctx, metric)
}

func GetAllMetrics(ctx context.Context) ([]metrics.Metric, error) {
	return serverStorage.GetAllMetrics(ctx)
}

func PingStorage(ctx context.Context) error {
	return serverStorage.PingStorage(ctx)
}

func GetAll(ctx context.Context) ([]string, error) {
	return serverStorage.GetAll(ctx)
}

func Close(ctx context.Context) {
	serverStorage.Close(ctx)
}
