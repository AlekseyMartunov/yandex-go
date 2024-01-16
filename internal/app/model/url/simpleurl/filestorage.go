// Package simpleurl uses when database dsn does not exist in config package
package simpleurl

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgconn"
)

// FileStorage uses when db dose not exist or connection unavailable
type FileStorage struct {
	filePath  string
	data      map[string]string
	currentID int
	sync.Mutex
}

// NewFileStorage return new struct
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

// Save new url pair
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

// Get return origin url
func (s *FileStorage) Get(key string) (string, error) {
	val, ok := s.data[key]

	if !ok {
		return "", ErrEmptyKey
	}
	return val, nil
}

// SaveBatch save several different urls
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

// GetShorted return shorted url by original
func (s *FileStorage) GetShorted(key string) (string, bool) {
	for k, v := range s.data {
		if v == key {
			return k, true
		}
	}
	return "", false
}

// GetAllURL return all url
func (s *FileStorage) GetAllURL(userID string) ([][2]string, error) {
	return nil, nil
}

// DeleteURL delete
func (s *FileStorage) DeleteURL(...URLToDel) error {
	return nil
}

// Ping just mocks for db interface
func (s *FileStorage) Ping() error {
	return errors.New("this is a file")
}

// Statistics return how urls and users store in application
func (s *FileStorage) Statistics() (int, int) {
	return 0, 0
}

// Close just mocks for db interface
func (s *FileStorage) Close() error {
	return nil
}
