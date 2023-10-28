package main

import (
	"context"
	"fmt"
	"log"
	_ "net/http/pprof"

	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/server/app"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/storagesync"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
)

var (
	Version string = "N/A"
	Date    string = "N/A"
	Commit  string = "N/A"
)

func printBuildInfo() {
	fmt.Printf("Build version: %s\n", Version)
	fmt.Printf("Build date: %s\n", Date)
	fmt.Printf("Build commit: %s\n", Commit)
}

func main() {

	printBuildInfo()

	cfg, err := config.LoadServerConfig()
	if err != nil {
		log.Fatalf("cannot load config: %s", err)
	}

	if err := logging.InitLogger(); err != nil {
		log.Fatalf("cannot init loger: %s", err)
	}

	storage.Init(context.Background(), cfg)

	if cfg.DatabaseAddress == "" {
		storagesync.RunSync(cfg)
	}

	app.Run(cfg)
	storage.Close(context.Background())
	if cfg.DatabaseAddress == "" {
		if err := storagesync.ShutdownSync(); err != nil {
			logging.Fatalf("cannot shutdown storage syncer: %s", err)
		}
	}

}
