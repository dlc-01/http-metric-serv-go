package storagesync

import (
	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var conf *config.ServerConfig
var shouldDumpMetricsOnMetrics bool
var workWithDB bool

var syncStor storage.Storage

func GetSyncMiddleware(databse string) gin.HandlerFunc {
	return func(gin *gin.Context) {
		if databse == "" {
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

func RunSync(cfg *config.ServerConfig, s storage.Storage) {
	syncStor = s
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

func ShutdownSync() error {
	return dumpFile()
}
