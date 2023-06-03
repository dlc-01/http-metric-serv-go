package routine

import (
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/hashing"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"net/http"
)

func sendMetrics(cfg *config.AgentConfig, result chan []metrics.Metric) {

	headers := map[string]string{
		"Content-Type":     "application/json",
		"Content-Encoding": "gzip",
	}

	jsons, err := metrics.ToJSONs(<-result)
	if err != nil {
		logging.Errorf("cannot generate request body: %w", err)
	}

	if cfg.HashKey != "" {
		headers["HashSHA256"] = hashing.HashingDate(cfg.HashKey, jsons)
	}

	gzip, err := metrics.Gzipper(jsons)
	if err != nil {
		logging.Errorf("cannot gzip body: %w", err)
	}

	resp, err := client.R().SetHeaders(headers).
		SetBody(gzip).
		Post(fmt.Sprintf("http://%s/updates/", cfg.ServerAddress))
	if err != nil {
		logging.Errorf("cannot generate request body: %w", err)
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusAccepted {
		logging.Errorf("unexpected status reponse code: %v", resp.StatusCode())
	}
}
