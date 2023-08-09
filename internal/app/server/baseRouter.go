package server

import "github.com/go-chi/chi/v5"

type storage interface {
	Encode(url string) string
	Decode(shortURL string) (string, bool)
}

type config interface {
	GetShorterURL() string
}

type baseRouter struct {
	db  storage
	cfg config
}

func NewBaseRouter(db storage, cfg config) *baseRouter {
	return &baseRouter{
		db:  db,
		cfg: cfg,
	}
}

func (br *baseRouter) Route() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/{url_id}", br.decodeURL)
	router.Post("/", br.encodeURL)

	return router
}
