package routine

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/dlc-01/http-metric-serv-go/internal/general/encryption"

	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/hashing"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
)

func sendMetrics(cfg *config.AgentConfig, metricsC chan []metrics.Metric) {

	wg := sync.WaitGroup{}
	wg.Add(cfg.LimitM)
	for i := 0; i < cfg.LimitM; i++ {
		go sendMetricsRoutine(&wg, metricsC, cfg)
	}
	wg.Wait()
}

func sendMetricsRoutine(wg *sync.WaitGroup, metricsC chan []metrics.Metric, cfg *config.AgentConfig) {
	headers := map[string]string{
		"Content-Type":     "application/json",
		"Content-Encoding": "gzip",
	}
	for items := range metricsC {
		jsons, err := metrics.ToJSON(items)
		if err != nil {
			logging.Errorf("cannot generate request body: %s", err)
			return
		}

		if cfg.PathCryptoKey != "" {
			encryptbuf, err := encryption.MetEncryptor.Encrypt(jsons)
			if err != nil {
				logging.Errorf("cannot encrypt metrics: %s", err)
			}

			jsons = encryptbuf
		}

		if cfg.HashKey != "" {
			headers["HashSHA256"] = hashing.HashingData(cfg.HashKey, jsons)
		}

		gzip, err := metrics.Gzipper(jsons)
		if err != nil {
			logging.Errorf("cannot gzip body: %s", err)
			return
		}

		resp, err := client.R().SetHeaders(headers).
			SetBody(gzip).
			Post(fmt.Sprintf("http://%s/updates/", cfg.ServerAddress))
		if err != nil {
			logging.Errorf("cannot generate request body: %s", err)
			return
		}
		if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusAccepted {
			logging.Errorf("unexpected status response code: %v", resp.StatusCode())
			return
		}

	}
	wg.Done()

}
