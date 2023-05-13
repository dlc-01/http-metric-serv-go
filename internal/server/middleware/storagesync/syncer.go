package storagesync

import (
	"context"
	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/storagesync/database"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/storagesync/file"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var ctx context.Context
var conf *config.ServerConfig
var shouldDumpMetricsOnMetrics bool
var workWithDB bool

func GetSyncMiddleware() gin.HandlerFunc {
	return func(gin *gin.Context) {
		gin.Next()
		if shouldDumpMetricsOnMetrics {
			if workWithDB {
				if err := database.DumpDB(); err != nil {
					logging.Warnf("cannot dump metrics to db %w", err)
				}
			} else {
				if err := file.DumpFile(); err != nil {
					logging.Fatalf("cannot dump file metrics to file: %s", err)
				}
			}

		}
		gin.Next()
	}
}

func RunSync(cfg *config.ServerConfig) {
	conf = cfg
	conf.DatabaseAddress = "postgresql://localhost:5432"
	conf.StoreInterval = 0
	ctx = context.Background()
	checkDB(conf.DatabaseAddress)
	if workWithDB {
		if err := database.InitDB(ctx, conf); err != nil {
			logging.Errorf("cannot init db: %s", err)
		}
	} else {
		file.InitForFile(conf)
	}
	if cfg.Restore {
		if workWithDB {
			if err := database.RestoreDB(); err != nil {
				logging.Warnf("cannot load metrics from db %w", err)
			}
		} else {
			if err := file.RestoreFile(); err != nil {
				logging.Warnf("cannot load metrics from file %w", err)
			}
		}

	}
	if conf.StoreInterval > 0 {
		if workWithDB {
			go database.RunDumperDB()
		} else {
			go file.RunDumperFile()
		}

	} else {
		shouldDumpMetricsOnMetrics = true
	}

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

func ShutdownSync() error {
	if workWithDB {
		return database.DumpDB()
	}
	return file.DumpFile()
}
