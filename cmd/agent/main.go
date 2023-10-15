package main

import (
	"context"
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/dlc-01/http-metric-serv-go/internal/agent/routine"
	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
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

	cfg, err := config.LoadAgentConfig()
	if err != nil {
		log.Fatalf("cannot load config: %s", err)
	}

	if err := logging.InitLogger(); err != nil {
		log.Fatalf("cannot init loger: %s", err)
	}

	storage.Init(context.Background(), &config.ServerConfig{})

	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go routine.Run(ctx, cfg)
	defer routine.Shutdown()
	logging.Info("agent has been started")

	<-term

	logging.Info("agent has been stopped")
}
