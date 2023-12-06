// Package encoder use to create unique short url
package encoder

// symbolsRunes, random key contain this runes
var symbolsRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// size is a len of random key
const size = 15

// generateRandomID create random id
func (e *Encoder) generateRandomID() string {

	id := e.getRandomValues()

	_, err := e.storage.Get(id)

	for err == nil {
		id = e.getRandomValues()
		_, err = e.storage.Get(id)
	}

	return id
}

// getRandomValues creates random values
func (e *Encoder) getRandomValues() string {

	output := make([]rune, size)
	for i := range output {
		output[i] = symbolsRunes[e.random.Intn(len(symbolsRunes))]
	}
	return string(output)
}
