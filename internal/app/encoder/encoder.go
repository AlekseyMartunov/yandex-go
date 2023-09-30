package encoder

import (
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/url/simpleurl"
)

type storage interface {
	Save(key, val, userID string) error
	Get(string) (string, error)
	SaveBatch(data *[][3]string, userID string) error
	GetShorted(key string) (string, bool)
	GetAllURL(userID string) ([][2]string, error)
	DeleteURL(...simpleurl.URLToDel) error
	Ping() error
}

type Encoder struct {
	storage storage
}

func NewEncoder(s storage) *Encoder {
	return &Encoder{storage: s}
}

func (e *Encoder) Encode(url, userID string) (string, error) {
	id := e.generateRandomID()
	err := e.storage.Save(id, url, userID)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (e *Encoder) Decode(id string) (string, error) {
	url, err := e.storage.Get(id)
	return url, err
}

func (e *Encoder) BatchEncode(data *[][3]string, userID string) error {
	// [[a, b, c], [a, b, c], ...]
	// a - CorrelationID
	// b - OriginalURL
	// c - ShortedURL

	for id := range *data {
		(*data)[id][2] = e.generateRandomID()
	}

	err := e.storage.SaveBatch(data, userID)
	if err != nil {
		return err
	}

	return nil
}

func (e *Encoder) GetShorted(url string) (string, bool) {
	shorted, ok := e.storage.GetShorted(url)
	return shorted, ok
}

func (e *Encoder) GetAllURL(userID string) ([][2]string, error) {
	res, err := e.storage.GetAllURL(userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (e *Encoder) DeleteURL(messages ...simpleurl.URLToDel) error {
	return e.storage.DeleteURL(messages...)
}

func (e *Encoder) Ping() error {
	return e.storage.Ping()
}
