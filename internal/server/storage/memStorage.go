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

var ms MemStorage

func Init() {
	ms.gauges = make(map[string]float64)
	ms.counters = make(map[string]int64)
}
func SetGauge(k string, v float64) {
	ms.gauges[k] = v
}

func GetGauge(k string) (float64, bool) {
	v, exist := ms.gauges[k]
	return v, exist
}

func SetCounter(k string, v int64) {
	if _, ok := ms.counters[k]; !ok {
		ms.counters[k] = 0
	}

	ms.counters[k] += v
}

func GetCounter(k string) (int64, bool) {
	v, exist := ms.counters[k]
	return v, exist
}
func GetAll() []string {
	names := make([]string, 0)
	for cm := range ms.counters {
		names = append(names, cm)
	}
	for gm := range ms.gauges {
		names = append(names, gm)
	}
	return names
}
