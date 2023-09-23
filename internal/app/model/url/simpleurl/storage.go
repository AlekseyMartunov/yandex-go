package simpleurl

type Storage interface {
	Save(key, val, userID string) error
	Get(key string) (string, bool)
	SaveBatch(data *[][3]string, userID string) error
	GetShorted(key string) (string, bool)
	GetAllURL(userID string) ([][2]string, error)
	DeleteURL(...URLToDel) error
	Ping() error
	Close() error
}

type fileLine struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type URLToDel struct {
	UserId string
	URL    string
}
