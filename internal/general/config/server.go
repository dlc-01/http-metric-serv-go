package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ServerConfig — structure for starting and running the server.
type ServerConfig struct {
	ServerAddress   string `json:"address"`        // server startup address
	StoreInterval   int    `json:"store_interval"` // data storage interval
	FileStoragePath string `json:"store_file"`     // storage path
	Restore         bool   `json:"restore"`        // data restoring check
	DatabaseAddress string `json:"database_dsn"`   // database address
	HashKey         string `json:"hash_key"`       // hash key
	LimitM          int    `json:"limit_m"`        // limit to receive metric
	PathCryptoKey   string `json:"crypto_key"`     // path for cryptoKey
	Config          string //  path to config in JSON
}

// LoadServerConfig — function to load data for server startup by
// means of flags and environment variables.
func LoadServerConfig() (*ServerConfig, error) {
	cfg := &ServerConfig{}
	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "server address")
	flag.IntVar(&cfg.StoreInterval, "i", 300, "store time interval")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/runtimeMetrics-db.json", "file data path")
	flag.BoolVar(&cfg.Restore, "r", true, "restore data")
	flag.StringVar(&cfg.DatabaseAddress, "d", "", "database address")
	flag.StringVar(&cfg.HashKey, "k", "", "hash key")
	flag.IntVar(&cfg.LimitM, "l", 8, "limit to receive metric")
	flag.StringVar(&cfg.PathCryptoKey, "crypto-key", "", "path to public crypto key")
	flag.StringVar(&cfg.Config, "c", "", "path to config in json")
	flag.StringVar(&cfg.Config, "config", "", "path to config in json")
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
	if envPathCryptoKey := os.Getenv("CRYPTO_KEY"); envPathCryptoKey != "" {
		cfg.PathCryptoKey = envPathCryptoKey
	}
	if envPathConfig := os.Getenv("CONFIG"); envPathConfig != "" {
		cfg.Config = envPathConfig
	}
	if cfg.Config != "" {
		var err error
		cfg, err = configFromJsonServer(cfg)
		if err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func configFromJsonServer(cfg *ServerConfig) (*ServerConfig, error) {
	f, err := os.ReadFile(cfg.Config)
	if err != nil {
		return nil, fmt.Errorf("can`t open file: %w", err)
	}
	newCfg := map[string]interface{}{}
	err = json.Unmarshal(f, &newCfg)
	if err != nil {
		return nil, fmt.Errorf("can`t unmarshal json: %w", err)
	}
	for key, value := range newCfg {
		switch key {
		case "address":
			cfg.ServerAddress = value.(string)
		case "store_interval":
			if strings.HasSuffix(value.(string), "s") {
				value, _ = strings.CutSuffix(value.(string), "s")
			}
			if storeInt, err := strconv.Atoi(value.(string)); err == nil {
				cfg.StoreInterval = storeInt

			} else {
				return nil, fmt.Errorf("cannot convert STORE_INTERVAL to int: %w", err)
			}
		case "store_file":
			cfg.FileStoragePath = value.(string)
		case "restore":
			cfg.Restore = value.(bool)
		case "database_dsn":
			cfg.DatabaseAddress = value.(string)
		case "hash_key":
			cfg.HashKey = value.(string)
		case "limit_m":
			cfg.LimitM = int(value.(float64))
		case "crypto_key":
			cfg.PathCryptoKey = value.(string)
		}
	}
	return cfg, nil
}
