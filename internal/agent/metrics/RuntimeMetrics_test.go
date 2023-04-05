package metrics

import (
	"reflect"
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
				"Alloc":  5.454,
				"Pizdec": 5.7,
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

func TestMemMetrics_GenerateURLMetrics(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	type args struct {
		host string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "first",
			fields: fields{gauge: map[string]float64{
				"alloc": 665.3,
			}},
			args: args{host: "host:8080"},
			want: "host:8080/update/gauge/alloc/665.3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics := &MemMetrics{
				gauge:   tt.fields.gauge,
				counter: tt.fields.counter,
			}
			if got := metrics.GenerateURLMetrics(tt.args.host); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateURLMetrics() = %v, want %v", got, tt.want)
			}
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
