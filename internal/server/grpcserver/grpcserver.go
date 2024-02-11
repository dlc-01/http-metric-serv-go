package grpcserver

import (
	"context"
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	pb "github.com/dlc-01/http-metric-serv-go/internal/protobuf"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
			err = status.Error(codes.NotFound, "can't find any metrics")
			return nil, err
		} else {
			if err = storage.SetMetric(ctx, metrics.Metric{ID: in.Metric.Name, MType: in.Metric.Type, Value: &in.Metric.Gauge}); err != nil {
				err = status.Errorf(codes.Internal, fmt.Sprintf("cannot set metric: %v", err))
				logging.Errorf("%w", err)
				return nil, err
			}
			response.Metric = in.Metric
		}

	case metrics.CounterType:
		if in.Metric.Counter < 0 {
			err = status.Error(codes.NotFound, "can't find any metrics")
		} else {
			if err = storage.SetMetric(ctx, metrics.Metric{ID: in.Metric.Name, MType: in.Metric.Type, Delta: &in.Metric.Counter}); err != nil {
				err = status.Errorf(codes.Internal, fmt.Sprintf("cannot set metric: %v", err))
				logging.Errorf("%w", err)
				return nil, err
			}
			response.Metric = in.Metric

		}

	default:
		response.Metric = nil
		err = status.Error(codes.NotFound, "cannot found that type metric")
		logging.Errorf("error: %w", err)
		return nil, err
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
				err = status.Error(codes.NotFound, "can't find any metrics")
				return nil, err
			} else {
				metric.Value = &m.Gauge
			}

		case metrics.CounterType:
			if m.Counter < 0 {
				err = status.Error(codes.NotFound, "can't find any metrics")
				return nil, err
			} else {
				metric.Delta = &m.Counter

			}

		default:
			err = status.Error(codes.NotFound, "cannot found that type metric")
			logging.Errorf("error: %w", err)
			return nil, err
		}
		if response.Success {
			metricRequest = append(metricRequest, metric)
		}
	}
	if err = storage.SetMetricsBatch(ctx, metricRequest); err != nil {
		err = status.Errorf(codes.Internal, fmt.Sprintf("error while setting metric butch : %v", err))
		logging.Errorf("%w", err)
		return nil, err
	}
	return &response, err
}
