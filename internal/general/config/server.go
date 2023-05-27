package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		if strings.HasSuffix(envStoreInterval, "s") {

			envStoreInterval, _ = strings.CutSuffix(envStoreInterval, "s")
			//need implements some methods while try catch error
			//return nil , fmt.Errorf("cannot cut s from STORE_INTERVAL: %w", err)
		}
		if storeInt, err := strconv.Atoi(envStoreInterval); err == nil {
			cfg.StoreInterval = storeInt

		} else {
			return nil, fmt.Errorf("cannot convert STORE_INTERVAL to int: %w", err)
		}
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		cfg.FileStoragePath = envFileStoragePath

	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		if restoreBoll, err := strconv.ParseBool(envRestore); err == nil {
			cfg.Restore = restoreBoll

		} else {
			return nil, fmt.Errorf("cannot convert RESTORE to boolean: %w", err)
		}
	}
	if envDatabasePath := os.Getenv("DATABASE_DSN"); envDatabasePath != "" {
		cfg.DatabaseAddress = envDatabasePath

	}
	return cfg, nil
}
