package storage

import (
	"math/rand"
)

var symbolsRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type storage struct {
	storage map[string]string
}

func NewStorage() *storage {
	s := storage{}
	s.storage = make(map[string]string)
	return &s
}

func (s *storage) Encode(url string) string {
	id := generateRandomID(10)
	_, ok := s.storage[id]

	for ok {
		id := generateRandomID(10)
		_, ok = s.storage[id]
	}

	s.storage[id] = url
	return id
}

func (s *storage) Decode(id string) (string, bool) {
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
