package storage

import (
	"bufio"
	"encoding/json"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/logging"
	"os"
)

func Restore(filename string) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logging.SLog.Info(err)
	}
	reader := bufio.NewReader(file)

	data, err := reader.ReadBytes('\n')

	if err != nil {
		logging.SLog.Info(err)
	}

	if err := json.Unmarshal(data, &ms); err != nil {
		logging.SLog.Info(err)
	}
}

func Save(filename string) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logging.SLog.Info(err)
	}

	writer := bufio.NewWriter(file)

	data, err := json.Marshal(&ms)
	if err != nil {
		logging.SLog.Info(err)
	}

	if _, err := writer.Write(data); err != nil {
		logging.SLog.Info(err)
	}

	if err := writer.WriteByte('\n'); err != nil {
		logging.SLog.Info(err)
	}

	if err := writer.Flush(); err != nil {
		logging.SLog.Info(err)
	}
}
