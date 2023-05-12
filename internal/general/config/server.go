package config

import (
	"flag"
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"os"
	"strconv"
)

type ServerConfig struct {
	ServerAddress   string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
	DatabaseAddress string
}

func LoadServerConfig() (*ServerConfig, error) {
	cfg := &ServerConfig{}
	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "server address")
	flag.IntVar(&cfg.StoreInterval, "i", 300, "store time interval")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/runtimeMetrics-db.json", "file data path")
	flag.BoolVar(&cfg.Restore, "r", true, "restore data")
	flag.StringVar(&cfg.DatabaseAddress, "d", "", "database address")
	flag.Parse()
	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		cfg.ServerAddress = envServerAddress
		logging.Infof("ADRESS: %s", envServerAddress)
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		if storeInt, err := strconv.Atoi(envStoreInterval); err == nil {
			cfg.StoreInterval = storeInt
			logging.Infof("STORE_INTERVAL: %s", storeInt)
		} else {
			return nil, fmt.Errorf("cannot convert STORE_INTERVAL to int: %w", err)
		}
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		cfg.FileStoragePath = envFileStoragePath
		logging.Infof("FILE_STORAGE_PATH: %s", envFileStoragePath)
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		if restoreBoll, err := strconv.ParseBool(envRestore); err == nil {
			cfg.Restore = restoreBoll
			logging.Infof("RESTORE: %s", restoreBoll)
		} else {
			return nil, fmt.Errorf("cannot convert RESTORE to boolean: %w", err)
		}
	}
	if envDatabasePath := os.Getenv("DATABASE_DSN"); envDatabasePath != "" {
		cfg.DatabaseAddress = envDatabasePath
		logging.Infof("DATABASE_DSN: %s", envDatabasePath)
	}
	return cfg, nil
}
