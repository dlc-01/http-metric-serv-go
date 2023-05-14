package routine

import (
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
)

func sendMetrics(addr string) error {
	jsons, err := metrics.ToJSONMetrics(storage.GetMetrics())
	if err != nil {
		return fmt.Errorf("cannot generate request body: %w", err)
	}
	//for _, m := range storage.GetMetrics() {
	//	json, err := m.ToJSON()
	//	if err != nil {
	//		return fmt.Errorf("cannot generate request body: %w", err)
	//	}
	//
	//	_, err = client.R().SetHeader("Content-Encoding", "gzip").
	//		SetHeader("Content-Type", "application/json").
	//		SetBody(json).
	//		Post(fmt.Sprintf("http://%s/update/", addr))
	//	if err != nil {
	//		return fmt.Errorf("cannot generate request body: %w", err)
	//	}
	//}
	_, err = client.R().SetHeader("Content-Encoding", "gzip").
		SetHeader("Content-Type", "application/json").
		SetBody(jsons).
		Post(fmt.Sprintf("http://%s/updates/", addr))
	if err != nil {
		return fmt.Errorf("cannot generate request body: %w", err)
	}
	return nil
}
