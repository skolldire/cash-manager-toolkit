package file_util

import (
	"os"
)

func ListFiles(dir string) ([]string, error) {
	var files []string
	read, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, file := range read {
		files = append(files, file.Name())
	}
	return files, nil
}
