package txt_file

import (
	"encoding/json"
	"io"
	"os"
)

func WriteStructsToFile[I any](path string, data []I) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, item := range data {
		jsonData, err := json.Marshal(item)
		if err != nil {
			return err
		}
		str := string(jsonData)
		str = str[1 : len(str)-1]
		str = str + "\n"
		_, err = file.WriteString(str)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadFile[O any](path string) ([]O, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var output []O
	err = json.Unmarshal(bytes, &output)
	if err != nil {
		return nil, err
	}
	return output, nil
}
