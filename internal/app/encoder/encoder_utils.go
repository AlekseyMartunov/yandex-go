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

func generateRandomID() string {
	output := make([]rune, size)
	for i := range output {
		output[i] = symbolsRunes[rand.Intn(len(symbolsRunes))]
	}
	return string(output)
}
