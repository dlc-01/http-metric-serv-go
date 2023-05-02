package params

import (
	"flag"
	"os"
	"strconv"
)

var (
	ServerAddress   string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
)

func ParseFlagsOs() {
	flag.StringVar(&ServerAddress, "a", "localhost:8080", "server address")
	flag.IntVar(&StoreInterval, "i", 300, "store time interval")
	flag.StringVar(&FileStoragePath, "f", "/tmp/metrics-db.json", "file data path")
	flag.BoolVar(&Restore, "r", true, "restore data")
	flag.Parse()
	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		ServerAddress = envServerAddress
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		if storeInt, err := strconv.Atoi(envStoreInterval); err == nil {
			StoreInterval = storeInt
		}
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		FileStoragePath = envFileStoragePath
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		if restoreBoll, err := strconv.ParseBool(envRestore); err == nil {
			Restore = restoreBoll
		}
	}
}
