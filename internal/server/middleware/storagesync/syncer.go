package storagesync

import (
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var conf *config.ServerConfig
var shouldDumpMetricsOnMetrics bool
var workWithDB bool

func GetSyncMiddleware() gin.HandlerFunc {
	return func(gin *gin.Context) {

		gin.Next()
		if shouldDumpMetricsOnMetrics {
			if err := dumpFile(); err != nil {
				logging.Fatalf("cannot dump file metrics to file: %s", err)
			}
		}
		gin.Next()
	}
}

func RunSync(cfg *config.ServerConfig) error {
	conf = cfg

	if err := restoreFile(); err != nil {
		return fmt.Errorf("cannot restore from file : %w", err)
	}

	if conf.StoreInterval > 0 {
		go runDumperFile()
	} else {
		shouldDumpMetricsOnMetrics = true
	}
	return nil

}

func ShutdownSync() error {
	return dumpFile()
}
