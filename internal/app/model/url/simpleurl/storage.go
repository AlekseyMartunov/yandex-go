// Package simpleurl uses when database dsn does not exist in config package
package simpleurl

import "errors"

// Storage interface to show witch methods struct need
type Storage interface {
	Save(key, val, userID string) error
	Get(key string) (string, error)
	SaveBatch(data *[][3]string, userID string) error
	GetShorted(key string) (string, bool)
	GetAllURL(userID string) ([][2]string, error)
	DeleteURL(...URLToDel) error
	Statistics() (int, int)
	Ping() error
	Close() error
}

// default url errors
var (
	ErrDeletedURL = errors.New("this URL is deleted")
	ErrEmptyKey   = errors.New("empty key")
)

type fileLine struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// URLToDel helps delete url
type URLToDel struct {
	UserID string
	URL    string
}
