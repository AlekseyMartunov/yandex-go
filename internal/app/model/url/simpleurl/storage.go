package simpleurl

import "errors"

type Storage interface {
	Save(key, val, userID string) error
	Get(key string) (string, error)
	SaveBatch(data *[][3]string, userID string) error
	GetShorted(key string) (string, bool)
	GetAllURL(userID string) ([][2]string, error)
	DeleteURL(...URLToDel) error
	Ping() error
	Close() error
}

var DeletedURLError = errors.New("this URL is deleted")
var EmptyKeyError = errors.New("empty key")

type fileLine struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type URLToDel struct {
	UserID string
	URL    string
}
