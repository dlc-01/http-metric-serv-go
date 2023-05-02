package storage

import (
	"bufio"
	"encoding/json"
	"github.com/dlc-01/http-metric-serv-go/internal/server/middleware/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/server/params"
	"os"
)

func Restore() error {

	file, err := os.OpenFile(params.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logging.SLog.Error(err, "OpenFile Restore")
	}
	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		logging.SLog.Error(err, "scanner Restore")
	}

	data := scanner.Bytes()

	err = json.Unmarshal(data, &ms)
	if err != nil {
		logging.SLog.Error(err, " unmarshal json Restore")
	}
	return nil
}
func Save() error {
	file, err := os.OpenFile(params.FileStoragePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logging.SLog.Error(err, "OpenFile Save")
	}
	writer := bufio.NewWriter(file)

	data, err := json.Marshal(&ms)
	if err != nil {
		logging.SLog.Error(err, "Marshal Save")
	}
	if _, err := writer.Write(data); err != nil {
		logging.SLog.Error(err, "Writer Save")
	}
	if err := writer.WriteByte('\n'); err != nil {
		logging.SLog.Error(err, "WriterByte Save")
	}
	if err := writer.Flush(); err != nil {
		logging.SLog.Error(err, "Flush Save")
	}
	return nil
}
