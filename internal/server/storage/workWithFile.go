package storage

import (
	"bufio"
	"encoding/json"
	"github.com/dlc-01/http-metric-serv-go/internal/server/params"
	"os"
)

func Restore() error {

	file, err := os.OpenFile(params.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil
	}

	data := scanner.Bytes()

	err = json.Unmarshal(data, &ms)
	if err != nil {
		return err
	}
	return nil
}
func Save() error {
	file, err := os.OpenFile(params.FileStoragePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)

	data, err := json.Marshal(&ms)
	if err != nil {
		return err
	}
	if _, err := writer.Write(data); err != nil {
		return err
	}
	if err := writer.WriteByte('\n'); err != nil {
		return err
	}
	if err := writer.Flush(); err != nil {
		return err
	}
	return nil
}
