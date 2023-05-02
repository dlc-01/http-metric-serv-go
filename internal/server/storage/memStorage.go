package storage

type Storage interface {
	setGauge(string, float64)
	setCounter(string, int64)
	getGauge(string) int64
	getCounter(string) float64
}

type MemStorage struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

var ms MemStorage

func Init() {
	ms.Gauges = make(map[string]float64)
	ms.Counters = make(map[string]int64)
}
func SetGauge(k string, v float64) {
	ms.Gauges[k] = v
}

func GetGauge(k string) (float64, bool) {
	v, exist := ms.Gauges[k]
	return v, exist
}

func SetCounter(k string, v int64) {
	if _, ok := ms.Counters[k]; !ok {
		ms.Counters[k] = 0
	}

	ms.Counters[k] += v
}

func GetCounter(k string) (int64, bool) {
	v, exist := ms.Counters[k]
	return v, exist
}
func GetAll() []string {
	names := make([]string, 0)
	for cm := range ms.Counters {
		names = append(names, cm)
	}
	for gm := range ms.Gauges {
		names = append(names, gm)
	}
	return names
}
