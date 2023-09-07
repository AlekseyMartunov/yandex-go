package simplestotage

import (
	"bufio"
	"encoding/json"
	"os"
)

type Storage interface {
	Save(key, val string) error
	Get(key string) (string, bool)
	SaveBatch(data *[][3]string) error
	Close() error
}

type fileLine struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewStorage(filePath string) Storage {

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return &MapStorage{data: make(map[string]string)}
	}
	defer file.Close()

	data := make(map[string]string)

	fl := &fileLine{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		json.Unmarshal(scanner.Bytes(), fl)
		data[fl.ShortURL] = fl.OriginalURL
	}
	s := FileStorage{
		filePath:  filePath,
		data:      data,
		currentID: fl.UUID + 1,
	}

	return &s
}
