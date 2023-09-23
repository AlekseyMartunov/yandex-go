package simpleurl

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/handlers"
	"github.com/jackc/pgx/v5/pgconn"
	"os"
	"sync"
)

type FileStorage struct {
	filePath  string
	data      map[string]string
	currentID int
	sync.Mutex
}

func NewFileStorage(filePath string) (Storage, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
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

	return &s, nil
}

func (s *FileStorage) Save(key, val, userID string) error {

	for _, v := range s.data {
		if v == val {
			return &pgconn.PgError{Code: "23505"}
		}
	}
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

func (s *FileStorage) SaveBatch(data *[][3]string, userID string) error {
	// [[a, b, c], [a, b, c], ...]
	// a - CorrelationID
	// b - OriginalURL
	// c - ShortedURL

	for id := range *data {
		key := (*data)[id][2]
		val := (*data)[id][1]
		err := s.Save(key, val, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *FileStorage) GetShorted(key string) (string, bool) {
	for k, v := range s.data {
		if v == key {
			return k, true
		}
	}
	return "", false
}

func (s *FileStorage) GetAllURL(userID string) ([][2]string, error) {
	return nil, nil
}

func (s *FileStorage) DeleteURL(...handlers.URLToDel) error {
	return nil
}

func (s *FileStorage) Ping() error {
	return errors.New("this is a file")
}

func (s *FileStorage) Close() error {
	return nil
}
