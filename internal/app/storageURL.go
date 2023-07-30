package app

import (
	"math/rand"
)

var symbolsRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

const host = "http://localhost:8080/"

func (s *app) encode(url string) string {
	id := generateRandomID(10)
	_, ok := s.storage[id]

	for ok {
		id := generateRandomID(10)
		_, ok = s.storage[id]
	}

	s.storage[id] = url
	return host + id
}

func (s *app) decode(id string) (string, bool) {
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
