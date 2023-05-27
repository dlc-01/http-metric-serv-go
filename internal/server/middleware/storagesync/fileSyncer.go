package storagesync

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/dlc-01/http-metric-serv-go/internal/server/handlers"
	"os"
	"time"
)

func runDumperFile() {
	dumpTicker := time.NewTicker(time.Duration(conf.StoreInterval) * time.Second)
	defer dumpTicker.Stop()
	for range dumpTicker.C {
		if err := dumpFile(); err != nil {
			logging.Fatalf("cannot dumpFile metrics to file: %s", err)
		}
	}
}

func restoreFile() error {

	file, err := os.OpenFile(conf.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return errors.New("cannot scan file")
	}

	data := scanner.Bytes()
	var new []metrics.Metric
	err = json.Unmarshal(data, &new)
	if err != nil {
		return fmt.Errorf("cannot decode line: %s", err)
	}
	err = syncStor.SetMetricsBatch(context.TODO(), new)
	if err != nil {
		return fmt.Errorf("cannot get storage: %w", err)
	}
	handlers.ServerStor.Storage = syncStor
	return nil
}

func dumpFile() error {
	syncStor = handlers.ServerStor.Storage
	file, err := os.OpenFile(conf.FileStoragePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	old, err := syncStor.GetAllMetrics(context.TODO())
	if err != nil {
		return fmt.Errorf("cannot get storage: %w", err)
	}
	data, err := json.Marshal(&old)
	if err != nil {
		return fmt.Errorf("cannot marshal metrics: %w", err)
	}
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("cannot encode runtimeMetrics: %w", err)
	}

	return nil
}
