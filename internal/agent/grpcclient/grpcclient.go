package grpcclient

import (
	"context"
	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	pb "github.com/dlc-01/http-metric-serv-go/internal/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
)

func SendMetricsViaGrpc(cfg *config.AgentConfig, metricsC chan []metrics.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(cfg.LimitM)
	conn, err := grpc.Dial(
		cfg.ServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logging.Errorf("error while connecting to GRPC server: %w", err)
	}

	c := pb.NewMetricsServiceClient(conn)

	for i := 0; i < cfg.LimitM; i++ {
		go sendBatchMetricsViaGrpc(&wg, c, metricsC)
	}
	wg.Wait()
}

func SendMetricViaGrpc(cfg *config.AgentConfig, metricsC chan metrics.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(cfg.LimitM)
	conn, err := grpc.Dial(
		cfg.ServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logging.Errorf("error while connecting to GRPC server: %w", err)
	}

	c := pb.NewMetricsServiceClient(conn)

	for i := 0; i < cfg.LimitM; i++ {
		go sendMetricViaGrpc(&wg, c, metricsC)
	}
	wg.Wait()
}

func sendMetricViaGrpc(wg *sync.WaitGroup, c pb.MetricsServiceClient, mCh <-chan metrics.Metric) {
	for m := range mCh {
		newM := castMetricToGrpc(m)
		_, err := c.UpdateMetric(context.Background(), &pb.UpdateMetricRequest{Metric: &newM})
		if err != nil {
			logging.Errorf("error while sending metric via grpc %w", err)
		}
	}
	wg.Done()

}

func sendBatchMetricsViaGrpc(wg *sync.WaitGroup, c pb.MetricsServiceClient, mCh <-chan []metrics.Metric) {
	ctx := context.Background()
	for m := range mCh {
		metric := castMetricsToGrpc(m)
		_, err := c.UpdateBatchMetrics(ctx, &pb.UpdateBatchMetricsRequest{Metrics: metric})
		if err != nil {
			logging.Errorf("error while sending metric via grpc %w", err)
		}
	}
	wg.Done()

}

func castMetricToGrpc(m metrics.Metric) pb.Metric {
	metric := pb.Metric{Name: m.ID, Type: m.MType}
	switch metric.Type {
	case metrics.GaugeType:
		metric.Gauge = *m.Value
	case metrics.CounterType:
		metric.Counter = *m.Delta
	}
	return metric
}

func castMetricsToGrpc(m []metrics.Metric) []*pb.Metric {
	metricGrpc := make([]*pb.Metric, 0)
	for _, v := range m {
		metric := castMetricToGrpc(v)
		metricGrpc = append(metricGrpc, &metric)
	}
	return metricGrpc
}
