package hashing

import (
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChekingHash(t *testing.T) {
	testValue := 2022.02
	testDelta := int64(24)
	logging.InitLogger()

	tests := []struct {
		name           string
		encodeKey      string
		decodeKey      string
		expectedResult bool
		responseBody   []metrics.Metric
	}{
		{
			name:           `true key`,
			expectedResult: true,
			encodeKey:      `secret_key`,
			decodeKey:      `secret_key`,
			responseBody: []metrics.Metric{
				{
					ID:    "TestCounter",
					MType: metrics.CounterType,
					Delta: &testDelta,
					Value: nil,
				},
				{
					ID:    "TestGauge",
					MType: metrics.GaugeType,
					Delta: nil,
					Value: &testValue,
				},
			},
		},
		{
			name:           `false key`,
			expectedResult: false,
			encodeKey:      `secret_key`,
			decodeKey:      `noSecret_key`,
			responseBody: []metrics.Metric{
				{
					ID:    "TestCounter",
					MType: metrics.CounterType,
					Delta: &testDelta,
					Value: nil,
				},
				{
					ID:    "TestGauge",
					MType: metrics.GaugeType,
					Delta: nil,
					Value: &testValue,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsons, _ := metrics.ToJSONs(tt.responseBody)
			hash := HashingDate(tt.decodeKey, jsons)
			bool, _ := CheckingHash(hash, tt.encodeKey, jsons)
			assert.Equal(t, tt.expectedResult, bool)
		})
	}
}
