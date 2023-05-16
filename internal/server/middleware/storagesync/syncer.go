package storagesync

import (
	"context"
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
			if workWithDB {
				if err := dumpDB(); err != nil {
					logging.Fatalf("cannot dump metrics to db %s", err)
				}
			} else {
				if err := dumpFile(); err != nil {
					logging.Fatalf("cannot dump file metrics to file: %s", err)
				}
			}
		}
		gin.Next()
	}
}

func RunSync(ctx context.Context, cfg *config.ServerConfig) {
	conf = cfg
	//conf.DatabaseAddress = "postgresql://localhost:5432"
	//conf.StoreInterval = 1
	db.ctx = ctx
	checkDB(conf.DatabaseAddress)
	if workWithDB {
		if err := initDB(); err != nil {
			logging.Errorf("cannot init db: %s", err)
		}
	}
	if conf.Restore {
		if workWithDB {
			if err := restoreDB(); err != nil {
				logging.Warnf("cannot load metrics from db %w", err)
			}
		} else {
			if err := restoreFile(); err != nil {
				logging.Warnf("cannot load metrics from file %w", err)
			}
		}

	}
	if conf.StoreInterval > 0 {
		if workWithDB {
			go runDumperDB()
		} else {
			go runDumperFile()
		}

	} else {
		shouldDumpMetricsOnMetrics = true
	}

}

func ShutdownSync() error {
	if workWithDB {
		return dumpDB()
	}
	return dumpFile()
}

func checkDB(DB string) {
	if DB != "" {
		workWithDB = true
		return
	} else {
		workWithDB = false
		return
	}
}
