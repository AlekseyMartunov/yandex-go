package encoder

type storage interface {
	Save(key, val string) error
	Get(string) (string, bool)
	SaveBatch(data *[][3]string) error
	GetShorted(key string) (string, bool)
}

type Encoder struct {
	storage storage
}

func NewEncoder(s storage) *Encoder {
	return &Encoder{storage: s}
}

func (e *Encoder) Encode(url string) (string, error) {
	id := e.generateRandomID()
	err := e.storage.Save(id, url)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (e *Encoder) Decode(id string) (string, bool) {
	url, ok := e.storage.Get(id)
	return url, ok
}

func (e *Encoder) BatchEncode(data *[][3]string) error {
	// [[a, b, c], [a, b, c], ...]
	// a - CorrelationID
	// b - OriginalURL
	// c - ShortedURL

	for id := range *data {
		(*data)[id][2] = e.generateRandomID()
	}

	err := e.storage.SaveBatch(data)
	if err != nil {
		return err
	}

	return nil
}

func (e *Encoder) GetShorted(url string) (string, bool) {
	shorted, ok := e.storage.GetShorted(url)
	return shorted, ok
}
