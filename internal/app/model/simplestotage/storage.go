package simplestotage

type Storage interface {
	Save(key, val string) error
	Get(key string) (string, bool)
	SaveBatch(data *[][3]string) error
	GetShorted(key string) (string, bool)
	Ping() error
	Close() error
}

type fileLine struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
