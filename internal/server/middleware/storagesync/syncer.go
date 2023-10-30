package storagesync

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
)

var conf *config.ServerConfig
var shouldDumpMetricsOnMetrics bool

// GetSyncMiddleware — middleware for gin that dump Metric to the file.
func GetSyncMiddleware() gin.HandlerFunc {
	return func(gin *gin.Context) {
		if conf.DatabaseAddress == "" {
			gin.Next()
			if shouldDumpMetricsOnMetrics {
				if err := dumpFile(); err != nil {
					logging.Fatalf("cannot dump file metrics to file: %s", err)
				}
			}
		}

		gin.Next()
	}
}

// RunSync — function that use for restoring data from the file and
// starting goroutine that dump data after a certain period of time.
func RunSync(cfg *config.ServerConfig) {

	conf = cfg
	if conf.Restore {
		if err := restoreFile(); err != nil {
			logging.Errorf("cannot restore from file : %s", err)
		}
	}
	if conf.StoreInterval > 0 {
		go runDumperFile()
	} else {
		shouldDumpMetricsOnMetrics = true
	}

}

// ShutdownSync  — function that use for dump data last type in the program.
func ShutdownSync() error {
	return dumpFile()
}
