package storage

import (
	"context"

	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
)

// Storage — interface for working with different storages.
type Storage interface {

	// Сreate — function that need for creating Storage.
	Сreate(ctx context.Context, cfg *config.ServerConfig) Storage

	// SetMetric — function that use for setting Metric.
	SetMetric(context.Context, metrics.Metric) error

	// SetMetricsBatch — function that setting BatchMetric.
	SetMetricsBatch(context.Context, []metrics.Metric) error

	// GetMetric — function that gets a metric by name and its type from Storage.
	GetMetric(context.Context, metrics.Metric) (metrics.Metric, error)

	// GetAllMetrics — function that gets all Metrics as an array of Metrics.
	GetAllMetrics(context.Context) ([]metrics.Metric, error)

	// GetAllStrings — function that gets all Metrics as an array of strings.
	GetAllStrings(context.Context) ([]string, error)

	// PingStorage — function that checks if the database is connected.
	PingStorage(context.Context) error

	// Close  — function that closing connection to db
	Close(context.Context)
}

// Init — function that initialize two types of storage (database and RAM).
func Init(ctx context.Context, conf *config.ServerConfig) {
	if conf.DatabaseAddress != "" {

		serverStorage = db.Сreate(ctx, conf)
		return
	}

	serverStorage = memS.Сreate(ctx, conf)
}

// serverStorage — variable that stores the Storage interface for further data proxying.
var serverStorage Storage

// SetMetric — proxying SetMetric in that package.
func SetMetric(ctx context.Context, metric metrics.Metric) error {
	return serverStorage.SetMetric(ctx, metric)
}

// SetMetricsBatch — proxying SetMetricsBatch in that package.
func SetMetricsBatch(ctx context.Context, metric []metrics.Metric) error {
	return serverStorage.SetMetricsBatch(ctx, metric)
}

// GetMetric — proxying GetMetric in that package.
func GetMetric(ctx context.Context, metric metrics.Metric) (metrics.Metric, error) {
	return serverStorage.GetMetric(ctx, metric)
}

// GetAllMetrics — proxying GetAllMetrics in that package.
func GetAllMetrics(ctx context.Context) ([]metrics.Metric, error) {
	return serverStorage.GetAllMetrics(ctx)
}

// PingStorage — proxying PingStorage in that package.
func PingStorage(ctx context.Context) error {
	return serverStorage.PingStorage(ctx)
}

// GetAll — proxying GetAll in that package.
func GetAll(ctx context.Context) ([]string, error) {
	return serverStorage.GetAllStrings(ctx)
}

// Close — proxying Close in that package.
func Close(ctx context.Context) {
	serverStorage.Close(ctx)
}
