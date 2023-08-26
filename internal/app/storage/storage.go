package storage

import (
	"bufio"
	"encoding/json"
	"os"
)

type fileLine struct {
	Uuid        int    `json:"uuid"`
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
	defer file.Close()

	if err != nil {
		panic(err)
	}

	fl := &fileLine{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		json.Unmarshal(scanner.Bytes(), fl)
	}

	return &Storage{
		currentID:       fl.Uuid + 1,
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
		Uuid:        s.currentID,
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
	defer file.Close()

	if err != nil {
		return err
	}

	_, err = file.Write(data)

	return err
}

func (s *Storage) Get(key string) (string, bool) {
	if !s.isFileAvailable {
		val, ok := s.data[key]
		return val, ok
	}

	file, err := os.OpenFile(s.filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer file.Close()

	if err != nil {
		panic(err)
	}

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
