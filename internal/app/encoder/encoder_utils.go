package encoder

import (
	"math/rand"
	"time"
)

var symbolsRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

const size = 15

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (e *Encoder) generateRandomID() string {

	id := getRandomValues()

	_, ok := e.storage.Get(id)

	for ok {
		id = getRandomValues()
		_, ok = e.storage.Get(id)
	}

	return id
}

func getRandomValues() string {
	output := make([]rune, size)
	for i := range output {
		output[i] = symbolsRunes[rand.Intn(len(symbolsRunes))]
	}
	return string(output)
}
