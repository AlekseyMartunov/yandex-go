package simplestotage

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"sync"
)

type MapStorage struct {
	data map[string]string
	sync.Mutex
}

func NewMapStorage() (Storage, error) {
	return &MapStorage{data: make(map[string]string)}, nil
}

func (s *MapStorage) Save(key, val string) error {
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

func (s *MapStorage) Get(key string) (string, bool) {
	val, ok := s.data[key]
	return val, ok
}

func (s *MapStorage) SaveBatch(data *[][3]string) error {
	// [[a, b, c], [a, b, c], ...]
	// a - CorrelationID
	// b - OriginalURL
	// c - ShortedURL

	for id := range *data {
		key := (*data)[id][2]
		val := (*data)[id][1]
		s.Save(key, val)
	}
	return nil
}

func (s *MapStorage) GetShorted(key string) (string, bool) {
	for k, v := range s.data {
		if v == key {
			return k, true
		}
	}
	return "", false
}

func (s *MapStorage) Ping() error {
	return errors.New("this is a map")
}

func (s *MapStorage) Close() error {
	return nil
}
