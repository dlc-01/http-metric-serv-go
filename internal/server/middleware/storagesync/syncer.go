package storagesync

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

var conf *config.ServerConfig
var shouldDumpMetricsOnMetrics bool

func GetSyncMiddleware(cfg *config.ServerConfig) gin.HandlerFunc {
	return func(gin *gin.Context) {
		logging.Infof("%+v/n", "GetSyncMiddleware")
		gin.Next()
		if shouldDumpMetricsOnMetrics {
			if err := dump(); err != nil {
				logging.Fatalf("cannot dump metrics to file: %s", err)
			}
		}
		gin.Next()
	}
}

func RunSync(cfg *config.ServerConfig) {
	conf = cfg
	if cfg.Restore {
		if err := restore(conf.FileStoragePath); err != nil {
			logging.Warnf("cannot load metrics from file %w", err)
		}
	}
	if conf.StoreInterval > 0 {
		go runDumper()
	} else {
		shouldDumpMetricsOnMetrics = true
	}

}

func ShutdownSync() error {
	return dump()
}

func runDumper() {
	dumpTicker := time.NewTicker(time.Duration(conf.StoreInterval) * time.Second)
	defer dumpTicker.Stop()
	for range dumpTicker.C {
		if err := dump(); err != nil {
			logging.Fatalf("cannot dump metrics to file: %s", err)
		}
	}
}

func restore(filePath string) error {

	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0666)
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

func dump() error {

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
