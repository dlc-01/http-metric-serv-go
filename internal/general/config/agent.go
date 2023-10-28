package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// AgentConfig — structure for starting and running agent of the server.
type AgentConfig struct {
	ServerAddress string `json:"address"`         // server client startup address
	Report        int    `json:"report_interval"` // report interval to server
	Poll          int    `json:"poll_interval"`   // poll interval to client
	HashKey       string `json:"hash_key"`        // hash key
	LimitM        int    `json:"limit_m"`         // limit to receive metric
	PathCryptoKey string `json:"crypto_key"`      // path for cryptoKey
	Config        string //  path to config in JSON
}

// LoadAgentConfig — function to load data for agent of server startup by
// means of flags and environment variables.
func LoadAgentConfig() (*AgentConfig, error) {
	cfg := &AgentConfig{}
	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "server address")
	flag.IntVar(&cfg.Report, "r", 10, "Report interval")
	flag.IntVar(&cfg.Poll, "p", 2, "Poll interval")
	flag.StringVar(&cfg.HashKey, "k", "", "hash key")
	flag.IntVar(&cfg.LimitM, "l", 8, "limit to collect metric")
	flag.StringVar(&cfg.PathCryptoKey, "crypto-key", "", "path to private crypto key")
	flag.StringVar(&cfg.Config, "config", "", "path to config in json")
	flag.StringVar(&cfg.Config, "c", "", "path to config in json")
	flag.Parse()

	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		cfg.ServerAddress = envServerAddress
	}
	if envReport := os.Getenv("REPORT_INTERVAL"); envReport != "" {

		if intReport, err := strconv.ParseInt(envReport, 10, 32); err == nil {
			cfg.Report = int(intReport)
		} else {
			return nil, fmt.Errorf("cannot parse REPORT_INTERVAL: %w", err)
		}

	}
	if envPoll := os.Getenv("POLL_INTERVAL"); envPoll != "" {
		if intPoll, err := strconv.ParseInt(envPoll, 10, 32); err == nil {
			cfg.Poll = int(intPoll)
		} else {
			return nil, fmt.Errorf("cannot parse POLL_INTERVAL: %w", err)
		}
	}
	if envHashKey := os.Getenv("KEY"); envHashKey != "" {
		cfg.HashKey = envHashKey
	}
	if envLimitM := os.Getenv("RATE_LIMIT"); envLimitM != "" {
		if intLimitM, err := strconv.ParseInt(envLimitM, 10, 32); err == nil {
			cfg.LimitM = int(intLimitM)
		} else {
			return nil, fmt.Errorf("cannot parse RATE_LIMIT: %w", err)
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
		cfg, err = configFromJSONAgent(cfg)
		if err != nil {
			return nil, err
		}
	}

	return cfg, nil

}

func configFromJSONAgent(cfg *AgentConfig) (*AgentConfig, error) {
	f, err := os.ReadFile(cfg.Config)
	if err != nil {
		return nil, fmt.Errorf("can`t open file: %w", err)
	}
	newCfg := map[string]interface{}{}
	err = json.Unmarshal(f, &newCfg)
	if err != nil {
		return nil, fmt.Errorf("can`t unmarshal json: %w", err)
	}
	for key, val := range newCfg {
		switch key {
		case "address":
			cfg.ServerAddress = val.(string)
		case "report_interval":
			if strings.HasSuffix(val.(string), "s") {
				val, _ = strings.CutSuffix(val.(string), "s")
			}
			if storeInt, err := strconv.Atoi(val.(string)); err == nil {
				cfg.Report = storeInt

			} else {
				return nil, fmt.Errorf("cannot convert report interval to int: %w", err)
			}
		case "pool_interval":
			if strings.HasSuffix(val.(string), "s") {
				val, _ = strings.CutSuffix(val.(string), "s")
			}
			if storeInt, err := strconv.Atoi(val.(string)); err == nil {
				cfg.Report = storeInt

			} else {
				return nil, fmt.Errorf("cannot convert pool interval to int: %w", err)
			}
		case "crypto_key":
			cfg.PathCryptoKey = val.(string)
		case "hash_key":
			cfg.HashKey = val.(string)
		case "limit_m":
			cfg.LimitM = val.(int)
		}
	}
	return cfg, nil
}
