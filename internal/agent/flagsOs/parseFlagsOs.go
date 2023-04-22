package flagsos

import (
	"flag"
	"log"
	"os"
	"strconv"
)

var (
	ServerAddress string
	Report        int
	Poll          int
)

func ParseFlagsOs() {
	flag.StringVar(&ServerAddress, "a", "localhost:8080", "server address")
	flag.IntVar(&Report, "r", 10, "Report interval")
	flag.IntVar(&Poll, "p", 2, "Poll interval")
	flag.Parse()

	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		ServerAddress = envServerAddress
	}

	if envReport := os.Getenv("REPORT_INTERVAL"); envReport != "" {
		intReport, err := strconv.ParseInt(envReport, 10, 32)
		if err != nil {
			log.Fatalf("cannot parse REPORT_INTERVAL: %v", err)
		}
		Report = int(intReport)
	}

	if envPoll := os.Getenv("POLL_INTERVAL"); envPoll != "" {
		intPoll, err := strconv.ParseInt(envPoll, 10, 32)
		if err != nil {
			log.Fatalf("cannot parse POLL_INTERVAL: %v", err)
		}
		Poll = int(intPoll)
	}
}
