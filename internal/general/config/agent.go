package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type AgentConfig struct {
	ServerAddress string
	Report        int
	Poll          int
	HashKey       string
}

func LoadAgentConfig() (*AgentConfig, error) {
	cfg := &AgentConfig{}
	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "server address")
	flag.IntVar(&cfg.Report, "r", 10, "Report interval")
	flag.IntVar(&cfg.Poll, "p", 2, "Poll interval")
	flag.StringVar(&cfg.HashKey, "k", "", "hash key")
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
	return cfg, nil
}
