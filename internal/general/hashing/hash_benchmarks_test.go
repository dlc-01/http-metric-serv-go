package hashing

import (
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/stretchr/testify/assert"
	"testing"
)

func BenchmarkCheckingHash(b *testing.B) {
	testValue := 2022.02
	testDelta := int64(24)
	logging.InitLogger()
	test := []metrics.Metric{{
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
		}}
	jsons, _ := metrics.ToJSON(test)
	hash := HashingData("key", jsons)
	b.StartTimer()
	bool, _ := CheckingHash(hash, "key", jsons)
	assert.Equal(b, true, bool)
}

func BenchmarkHashingDate(b *testing.B) {
	testValue := 2022.02
	testDelta := int64(24)
	logging.InitLogger()
	test := []metrics.Metric{{
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
		}}
	jsons, _ := metrics.ToJSON(test)
	b.StartTimer()
	hash := HashingData("key", jsons)
	b.StopTimer()
	bool, _ := CheckingHash(hash, "key", jsons)
	assert.Equal(b, true, bool)
}
