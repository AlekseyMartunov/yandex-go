package storage

import (
	"bufio"
	"encoding/json"
	"os"
)

type fileLine struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Storage struct {
	filePath        string
	currentID       int
	isFileAvailable bool
	data            map[string]string
}

func NewStorage(filePath string) (*Storage, error) {
	if filePath == "" {
		return &Storage{
			currentID:       1,
			data:            make(map[string]string),
			isFileAvailable: false,
		}, nil
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	fl := &fileLine{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		json.Unmarshal(scanner.Bytes(), fl)
	}

	return &Storage{
		currentID:       fl.UUID + 1,
		filePath:        filePath,
		isFileAvailable: true,
	}, nil
}

func (s *Storage) Save(key, val string) error {
	if !s.isFileAvailable {
		s.data[key] = val
		return nil
	}

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

func (s *Storage) Get(key string) (string, bool) {
	if !s.isFileAvailable {
		val, ok := s.data[key]
		return val, ok
	}

	file, err := os.OpenFile(s.filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fl := &fileLine{}
		json.Unmarshal(scanner.Bytes(), fl)

		if fl.ShortURL == key {
			return fl.OriginalURL, true
		}
	}

	return "", false
}
