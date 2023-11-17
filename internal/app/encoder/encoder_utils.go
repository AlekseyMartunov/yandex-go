package encoder

var symbolsRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

const size = 15

func (e *Encoder) generateRandomID() string {

	id := e.getRandomValues()

	_, err := e.storage.Get(id)

	for err == nil {
		id = e.getRandomValues()
		_, err = e.storage.Get(id)
	}

	return id
}

func (e *Encoder) getRandomValues() string {

	output := make([]rune, size)
	for i := range output {
		output[i] = symbolsRunes[e.random.Intn(len(symbolsRunes))]
	}
	return string(output)
}
