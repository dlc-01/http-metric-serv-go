package metrics

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
)

// Types of metric.
const (
	GaugeType   = "gauge"   // type gauge
	CounterType = "counter" // type counter
)

// Metric — structure for transferring the metric value to the server.
type Metric struct {
	ID    string   `json:"id"`              // name metric
	MType string   `json:"type"`            // parameter taking the value gauge or counter.
	Delta *int64   `json:"delta,omitempty"` //metric value in case of counter transfer.
	Value *float64 `json:"value,omitempty"` // metric value in case of gauge transfer.
}

// ToJSONWithGzip — function that formats the metric to json and then to gzip.
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

// ToJSON — function that formats an array of metrics into json.
func ToJSON(m []Metric) ([]byte, error) {
	json, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal request to json: %w", err)
	}
	return json, nil
}

// Gzipper — function that compresses json to gzip.
func Gzipper(json []byte) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	defer gz.Close()

	if _, err := gz.Write(json); err != nil {
		return nil, fmt.Errorf("cannot compresed data: %w", err)
	}
	return &buf, nil
}
