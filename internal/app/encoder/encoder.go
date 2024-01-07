// Package encoder use to create unique short url
package encoder

import (
	"math/rand"
	"time"

	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/url/simpleurl"
)

type storage interface {
	Save(key, val, userID string) error
	Get(string) (string, error)
	SaveBatch(data *[][3]string, userID string) error
	GetShorted(key string) (string, bool)
	GetAllURL(userID string) ([][2]string, error)
	DeleteURL(...simpleurl.URLToDel) error
	Statistics() (int, int)
	Ping() error
}

// Encoder type use to creates unique short url
type Encoder struct {
	random  *rand.Rand
	storage storage
}

// NewEncoder creates new struct
func NewEncoder(s storage) *Encoder {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &Encoder{
		storage: s,
		random:  r,
	}
}

// Encode creates new shorten url
func (e *Encoder) Encode(url, userID string) (string, error) {
	id := e.generateRandomID()
	err := e.storage.Save(id, url, userID)
	if err != nil {
		return "", err
	}
	return id, nil
}

// Decode return origin url
func (e *Encoder) Decode(id string) (string, error) {
	url, err := e.storage.Get(id)
	return url, err
}

// BatchEncode shortens several URLs at once
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

// GetShorted return shorted url by original
func (e *Encoder) GetShorted(url string) (string, bool) {
	shorted, ok := e.storage.GetShorted(url)
	return shorted, ok
}

// GetAllURL return all user s url
func (e *Encoder) GetAllURL(userID string) ([][2]string, error) {
	res, err := e.storage.GetAllURL(userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// DeleteURL deletes urp
func (e *Encoder) DeleteURL(messages ...simpleurl.URLToDel) error {
	return e.storage.DeleteURL(messages...)
}

func (e *Encoder) Statistics() (int, int) {
	return e.storage.Statistics()
}

// Ping checks store is available
func (e *Encoder) Ping() error {
	return e.storage.Ping()
}
