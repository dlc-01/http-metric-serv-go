package metrics

import (
	"testing"
)

func TestMemMetrics_Check(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	var tests = []struct {
		name   string
		fields fields
	}{
		{
			name: "firstTest",
			fields: fields{gauge: map[string]float64{
				"Alloc": 5.454,
				"Da":    5.7,
			}, counter: map[string]int64{
				"first":  1,
				"second": 2,
			},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics := &MemMetrics{
				gauge:   tt.fields.gauge,
				counter: tt.fields.counter,
			}
			metrics.Check()
		})
	}
}

func TestMemMetrics_Init(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "first",
			fields: fields{gauge: make(map[string]float64), counter: make(map[string]int64)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics := &MemMetrics{
				gauge:   tt.fields.gauge,
				counter: tt.fields.counter,
			}
			metrics.Init()
		})
	}
}
