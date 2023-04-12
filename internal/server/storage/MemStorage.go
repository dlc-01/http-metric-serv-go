package storage

type Storage interface {
	setGauge(string, float64)
	setCounter(string, int64)
	getGauge(string) int64
	getCounter(string) float64
}

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

var Ms MemStorage

func (Ms *MemStorage) Init() {
	(*Ms).gauges = make(map[string]float64)
	(*Ms).counters = make(map[string]int64)
}

func (Ms *MemStorage) SetGauge(k string, v float64) {
	Ms.gauges[k] = v
}

func (Ms *MemStorage) GetGauge(k string) (float64, bool) {
	v, exist := Ms.gauges[k]
	return v, exist
}

func (Ms *MemStorage) SetCounter(k string, v int64) {
	if _, ok := Ms.counters[k]; !ok {
		Ms.counters[k] = 0
	}

	Ms.counters[k] += v
}

func (Ms *MemStorage) GetCounter(k string) (int64, bool) {
	v, exist := Ms.counters[k]
	return v, exist
}
