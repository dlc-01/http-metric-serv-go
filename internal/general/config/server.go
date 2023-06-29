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
	HashKey         string
	LimitM          int
}

func LoadServerConfig() (*ServerConfig, error) {
	cfg := &ServerConfig{}
	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "server address")
	flag.IntVar(&cfg.StoreInterval, "i", 300, "store time interval")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/runtimeMetrics-db.json", "file data path")
	flag.BoolVar(&cfg.Restore, "r", true, "restore data")
	flag.StringVar(&cfg.DatabaseAddress, "d", "", "database address")
	flag.StringVar(&cfg.HashKey, "k", "", "hash key")
	flag.IntVar(&cfg.LimitM, "l", 8, "limit to receive metric")
	flag.Parse()
	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		cfg.ServerAddress = envServerAddress
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		if strings.HasSuffix(envStoreInterval, "s") {

			envStoreInterval, _ = strings.CutSuffix(envStoreInterval, "s")

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
		if restoreBool, err := strconv.ParseBool(envRestore); err == nil {
			cfg.Restore = restoreBool
		} else {
			return nil, fmt.Errorf("cannot convert RESTORE to boolean: %w", err)
		}
	}
	if envDatabasePath := os.Getenv("DATABASE_DSN"); envDatabasePath != "" {
		cfg.DatabaseAddress = envDatabasePath
	}
	if envHashKey := os.Getenv("KEY"); envHashKey != "" {
		cfg.HashKey = envHashKey
	}

	if envLimitM := os.Getenv("RATE_LIMIT"); envLimitM != "" {
		if intLimitM, err := strconv.ParseInt(envLimitM, 10, 32); err == nil {
			cfg.LimitM = int(intLimitM)
		} else {
			return nil, fmt.Errorf("cannot parseRATE_LIMIT: %w", err)
		}
	}
	return cfg, nil
}
