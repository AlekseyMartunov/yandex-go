package simplestotage

import (
	"encoding/json"
	"os"
	"sync"
)

type FileStorage struct {
	filePath  string
	data      map[string]string
	currentID int
	sync.Mutex
}

func (s *FileStorage) Save(key, val string) error {

	s.data[key] = val

	fl := fileLine{
		UUID:        s.currentID,
		ShortURL:    key,
		OriginalURL: val,
	}

	s.currentID++

	data, err := json.Marshal(fl)
	if err != nil {
		return err
	}

	data = append(data, '\n')

	file, err := os.OpenFile(s.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(data)

	return err
}

func (s *FileStorage) Get(key string) (string, bool) {
	val, ok := s.data[key]
	return val, ok
}

func (s *FileStorage) SaveBatch(data *[][3]string) error {
	// [[a, b, c], [a, b, c], ...]
	// a - CorrelationID
	// b - OriginalURL
	// c - ShortedURL

	for id, _ := range *data {
		key := (*data)[id][2]
		val := (*data)[id][1]
		err := s.Save(key, val)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *FileStorage) Close() error {
	return nil
}
