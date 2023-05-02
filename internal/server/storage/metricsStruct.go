package storage

type Metrics struct {
	ID    string   `jsonh:"id"`              // имя метрики
	MType string   `jsonh:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `jsonh:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `jsonh:"value,omitempty"` // значение метрики в случае передачи gauge
}
