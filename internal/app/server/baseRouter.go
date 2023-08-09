package router

import "github.com/go-chi/chi/v5"

type storage interface {
	Encode() string
	Decode() string
}

type config interface {
	GetAdres() string
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

func (r *baseRouter) Route() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/{url_id}", a.DecodeURL)
}
