package encoder

import (
	"math/rand"
)

var symbolsRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type storage interface {
	Save(string, string)
	Get(string) (string, bool)
}

type Encoder struct {
	storage storage
}

func NewEncoder(s storage) *Encoder {
	return &Encoder{storage: s}
}

func (e *Encoder) Encode(url string) string {
	id := generateRandomID(10)
	_, ok := e.storage.Get(id)

	for ok {
		id := generateRandomID(10)
		_, ok = e.storage.Get(id)
	}

	e.storage.Save(id, url)
	return id
}

func (e *Encoder) Decode(id string) (string, bool) {
	url, ok := e.storage.Get(id)
	return url, ok
}

func generateRandomID(size int) string {
	output := make([]rune, size)
	for i := range output {
		output[i] = symbolsRunes[rand.Intn(len(symbolsRunes))]
	}
	return string(output)
}
