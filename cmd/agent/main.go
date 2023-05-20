package main

import (
	"github.com/dlc-01/http-metric-serv-go/internal/agent/routine"
	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.LoadAgentConfig()
	if err != nil {
		log.Fatalf("cannot load config: %s", err)
	}

	if err := logging.InitLogger(); err != nil {
		log.Fatalf("cannot init loger: %s", err)
	}

	storage.Init()

	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGINT, syscall.SIGTERM)

	go routine.Run(cfg)

	logging.Info("agent has been started")

	<-term

	routine.Shutdown()

	logging.Info("agent has been stopped")
}
