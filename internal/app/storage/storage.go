package storage

import (
	"math/rand"
)

var symbolsRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type Storage struct {
	storage map[string]string
}

func NewStorage() *Storage {
	return &Storage{storage: make(map[string]string)}
}

func (s *Storage) Encode(url string) string {
	id := generateRandomID(10)
	_, ok := s.storage[id]

	for ok {
		id := generateRandomID(10)
		_, ok = s.storage[id]
	}

	s.storage[id] = url
	return id
}

func (s *Storage) Decode(id string) (string, bool) {
	url, ok := s.storage[id]
	return url, ok
}

func generateRandomID(size int) string {
	output := make([]rune, size)
	for i := range output {
		output[i] = symbolsRunes[rand.Intn(len(symbolsRunes))]
	}
	return string(output)
}
