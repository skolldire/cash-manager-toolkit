package file_util

import (
	"io"
	"os"
)

func ReadJsonFile(file string) []byte {
	jsonFile, err := os.Open(file)
	if err != nil {
		return nil
	}
	defer func(jsonFile *os.File) {
		_ = jsonFile.Close()
	}(jsonFile)
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil
	}
	return byteValue
}
