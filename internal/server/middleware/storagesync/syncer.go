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

func GetSyncMiddleware() gin.HandlerFunc {
	return func(gin *gin.Context) {

		gin.Next()

		if skipDumpMetrics {
			if err := dump(); err != nil {
				logging.Fatalf("cannot dump metrics to file: %s", err)
			}
		} else {
			//FIXME Called from everyone
			go runDumper()
		}

	}
}

var conf *config.ServerConfig
var skipDumpMetrics bool

func RunSync(cfg *config.ServerConfig) error {
	conf = cfg
	if cfg.Restore {
		if err := restore(conf.FileStoragePath); err != nil {
			return fmt.Errorf("cannot load metrics from file %w", err)
		}
	}
	if conf.StoreInterval > 0 {
		go runDumper()
	} else {
		skipDumpMetrics = true
	}

	return nil
}

func ShutdownSync() error {
	return dump()
}

func runDumper() {
	dumpTicker := time.NewTicker(time.Duration(conf.StoreInterval) * time.Second)

	for true {
		select {
		case <-dumpTicker.C:
			if err := dump(); err != nil {
				logging.Fatalf("cannot dump metrics to file: %s", err)
			}
		}
	}
}

func restore(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
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
	writer := bufio.NewWriter(file)

	old := storage.GetStorage()
	data, err := json.Marshal(&old)

	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	if _, err := writer.Write(data); err != nil {
		return fmt.Errorf("cannot encode runtimeMetrics: %w", err)
	}
	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("cannot write runtimeMetrics to the file: %w", err)
	}

	return nil
}
