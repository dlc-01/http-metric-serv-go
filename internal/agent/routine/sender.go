package routine

import (
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
)

func sendMetrics(addr string) error {
	for _, m := range storage.GetMetrics() {
		json, err := m.ToJSON()
		if err != nil {
			return fmt.Errorf("cannot generate request body: %w", err)
		}
		//FIXME ERROR: gzip: Invalid header if we will check error
		client.R().SetHeader("Content-Encoding", "gzip").
			SetHeader("Accept-Encoding", "gzip").
			SetHeader("Content-Type", "application/json").
			SetBody(json).
			Post(fmt.Sprintf("http://%s/update/", addr))

	}
	return nil
}
