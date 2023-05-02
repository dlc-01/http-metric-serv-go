package storage

import (
	"bufio"
	"encoding/json"
	"os"
)

func Restore(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(file)

	data, err := reader.ReadBytes('\n')
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &ms); err != nil {
		return err
	}

	return nil
}

func Save(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
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
