package simplestotage

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"
)

type Storage interface {
	Save(key, val string) error
	Get(key string) (string, bool)
	Close() error
}

type PostgresDB interface {
	Save(key, val string) error
	Get(key string) (string, bool)
}

type MapStorage struct {
	data map[string]string
	sync.Mutex
}

type FileStorage struct {
	filePath  string
	data      map[string]string
	currentID int
	sync.Mutex
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

func (s *MapStorage) Save(key, val string) error {
	s.data[key] = val
	return nil
}

func (s *MapStorage) Get(key string) (string, bool) {
	val, ok := s.data[key]
	return val, ok
}

func (s *MapStorage) Close() error {
	return nil
}

func (s *FileStorage) Close() error {
	return nil
}
