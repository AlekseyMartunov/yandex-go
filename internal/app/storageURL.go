package app

import (
	"math/rand"
)

var symbolsRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (a *app) encode(url string) string {
	id := generateRandomID(10)
	_, ok := a.storage[id]

	for ok {
		id := generateRandomID(10)
		_, ok = a.storage[id]
	}

	a.storage[id] = url
	return id
}

func (a *app) decode(id string) (string, bool) {
	url, ok := a.storage[id]
	return url, ok
}

func generateRandomID(size int) string {
	output := make([]rune, size)
	for i := range output {
		output[i] = symbolsRunes[rand.Intn(len(symbolsRunes))]
	}
	return string(output)
}
