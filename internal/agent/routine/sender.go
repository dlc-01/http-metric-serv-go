package routine

import (
	"context"
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"net/http"
)

func sendMetrics(addr string) error {
	metric, err := storage.ServerStorage.GetAllMetrics(context.TODO())
	if err != nil {
		return fmt.Errorf("cannot get metrics :%w", err)
	}
	jsons, err := metrics.ToJSONWithGzipMetrics(metric)
	if err != nil {
		return fmt.Errorf("cannot generate request body: %w", err)
	}
	resp, err := client.R().SetHeader("Content-Encoding", "gzip").
		SetHeader("Content-Type", "application/json").
		SetBody(jsons).
		Post(fmt.Sprintf("http://%s/updates/", addr))
	if err != nil {
		return fmt.Errorf("cannot generate request body: %w", err)
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusAccepted {
		return fmt.Errorf("unexpected status reponse code: %v", resp.StatusCode())
	}
	return nil
}
