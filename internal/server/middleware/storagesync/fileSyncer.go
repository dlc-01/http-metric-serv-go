package storagesync

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
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

	new := storage.GetStorage()
	err = json.Unmarshal(data, &new)
	if err != nil {
		return fmt.Errorf("cannot decode line: %s", err)
	}
	storage.SetStorage(new)

	return nil
}

func dumpFile() error {
	file, err := os.OpenFile(conf.FileStoragePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	old := storage.GetStorage()

	data, err := json.Marshal(&old)
	if err != nil {
		return fmt.Errorf("cannot marshal metrics: %w", err)
	}
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("cannot encode runtimeMetrics: %w", err)
	}
	if _, err := file.Write([]byte("\n")); err != nil {
		return fmt.Errorf("cannot write runtimeMetrics to the file: %w", err)
	}

	return nil
}
