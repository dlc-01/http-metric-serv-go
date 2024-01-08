package grpcserver

import (
	"context"
	"errors"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	pb "github.com/dlc-01/http-metric-serv-go/internal/protobuf"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
)

type GrpsServer struct {
	pb.UnimplementedMetricsServiceServer
}

func CreateGrpcServer() *GrpsServer {
	return &GrpsServer{}
}

func (s *GrpsServer) UpdateMetric(ctx context.Context, in *pb.UpdateMetricRequest) (*pb.UpdateMetricResponse, error) {
	var err error
	var response pb.UpdateMetricResponse
	switch in.Metric.Type {
	case metrics.GaugeType:
		if in.Metric.Gauge < 0 {
			response.Success = false
		} else {
			if err = storage.SetMetric(ctx, metrics.Metric{ID: in.Metric.Name, MType: in.Metric.Type, Value: &in.Metric.Gauge}); err != nil {
				logging.Errorf("cannot set metric: %w", err)
				response.Metric = nil
				response.Success = false
			}
			response.Metric = in.Metric
			response.Success = true
		}

	case metrics.CounterType:
		if in.Metric.Counter < 0 {
			response.Success = false
		} else {
			if err = storage.SetMetric(ctx, metrics.Metric{ID: in.Metric.Name, MType: in.Metric.Type, Delta: &in.Metric.Counter}); err != nil {
				logging.Errorf("cannot set metric: %w", err)
				response.Metric = nil
				response.Success = false
			}
			response.Metric = in.Metric
			response.Success = true
		}

	default:
		response.Metric = nil
		response.Success = false
		err = errors.New("can't find that type")
		logging.Errorf("error: %w", err)
	}
	return &response, err
}

func (s *GrpsServer) UpdateBatchMetrics(ctx context.Context, in *pb.UpdateBatchMetricsRequest) (*pb.UpdateBatchMetricsResponse, error) {
	var err error
	var response pb.UpdateBatchMetricsResponse
	metricRequest := make([]metrics.Metric, 0)
	for _, m := range in.Metrics {
		metric := metrics.Metric{ID: m.Name, MType: m.Type}
		switch m.Type {
		case metrics.GaugeType:
			if m.Gauge < 0 {
				response.Success = false
			} else {
				metric.Value = &m.Gauge
				response.Success = true
			}

		case metrics.CounterType:
			if m.Counter < 0 {
				response.Success = false
			} else {
				metric.Delta = &m.Counter

				response.Success = true
			}

		default:

			response.Success = false
			err = errors.New("can't find that type")
			logging.Errorf("error: %w", err)
		}
		if response.Success {
			metricRequest = append(metricRequest, metric)
		}
	}
	if err = storage.SetMetricsBatch(ctx, metricRequest); err != nil {
		logging.Errorf("error while setting metric butch : %w", err)
	}
	return &response, err
}
