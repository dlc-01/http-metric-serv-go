package metrics

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
)

const (
	GaugeType   = "gauge"
	CounterType = "counter"
)

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m *Metric) ToJSONWithGzip() (*bytes.Buffer, error) {
	json, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal request to json: %w", err)
	}
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	defer gz.Close()

	if _, err = gz.Write(json); err != nil {
		return nil, fmt.Errorf("cannot compresed data: %w", err)
	}
	return &buf, nil
}

func ToJSONWithGzipMetrics(m []Metric) (*bytes.Buffer, error) {
	json, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal request to json: %w", err)
	}
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	defer gz.Close()

	if _, err = gz.Write(json); err != nil {
		return nil, fmt.Errorf("cannot compresed data: %w", err)
	}
	return &buf, nil
}
