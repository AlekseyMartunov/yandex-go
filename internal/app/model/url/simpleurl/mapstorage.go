// Package simpleurl uses when database dsn does not exist in config package
package simpleurl

import (
	"errors"
	"sync"

	"github.com/jackc/pgx/v5/pgconn"
)

// MapStorage contains a map for store url
type MapStorage struct {
	data map[string]string
	sync.Mutex
}

// NewMapStorage returns new struct
func NewMapStorage() (Storage, error) {
	return &MapStorage{data: make(map[string]string)}, nil
}

// Save uses to save new url pair
func (s *MapStorage) Save(key, val, userID string) error {
	for _, v := range s.data {
		if v == val {
			return &pgconn.PgError{Code: "23505"}
		}
	}

	s.Mutex.Lock()
	s.data[key] = val
	s.Mutex.Unlock()
	return nil
}

// Get uses to get url
func (s *MapStorage) Get(key string) (string, error) {
	val, ok := s.data[key]

	if !ok {
		return "", ErrEmptyKey
	}

	return val, nil
}

// SaveBatch save several different urls
func (s *MapStorage) SaveBatch(data *[][3]string, userID string) error {
	// [[a, b, c], [a, b, c], ...]
	// a - CorrelationID
	// b - OriginalURL
	// c - ShortedURL

	for id := range *data {
		key := (*data)[id][2]
		val := (*data)[id][1]
		s.Save(key, val, userID)
	}
	return nil
}

// GetShorted return shorted url
func (s *MapStorage) GetShorted(key string) (string, bool) {
	for k, v := range s.data {
		if v == key {
			return k, true
		}
	}
	return "", false
}

// GetAllURL just mocks for db interface
func (s *MapStorage) GetAllURL(userID string) ([][2]string, error) {
	return nil, nil
}

// DeleteURL just mocks for db interface
func (s *MapStorage) DeleteURL(...URLToDel) error {
	return nil
}

// Ping ust mocks for db interface
func (s *MapStorage) Ping() error {
	return errors.New("this is a map")
}

// Statistics return how urls and users store in application
func (s *MapStorage) Statistics() (int, int) {
	return 0, 0
}

// Close just mocks for db interface
func (s *MapStorage) Close() error {
	return nil
}
