package main

import (
	"context"
	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/server/app"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/storagesync"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"log"
)

func main() {
	cfg, err := config.LoadServerConfig()
	if err != nil {
		log.Fatalf("cannot load config: %s", err)
	}

	if err := logging.InitLogger(); err != nil {
		log.Fatalf("cannot init loger: %s", err)
	}
	//cfg.DatabaseAddress = "postgresql://localhost:5432"
	storage.Init(context.Background(), cfg)

	if err := storagesync.RunSync(cfg); err != nil {
		logging.Warnf("cannot run syncer %s", err)
	}

	app.Run(cfg.ServerAddress)

	if err := storagesync.ShutdownSync(); err != nil {
		logging.Fatalf("cannot shutdown storage syncer: %s", err)
	}

}
