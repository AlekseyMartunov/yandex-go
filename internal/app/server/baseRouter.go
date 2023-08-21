package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ShortURLHandler interface {
	EncodeURL(w http.ResponseWriter, r *http.Request)
	DecodeURL(w http.ResponseWriter, r *http.Request)
}

type Logger interface {
	WithLogging(next http.HandlerFunc) http.HandlerFunc
}

type BaseRouter struct {
	handler ShortURLHandler
	logger  Logger
}

func NewBaseRouter(h ShortURLHandler, l Logger) *BaseRouter {
	return &BaseRouter{handler: h, logger: l}
}

func (br *BaseRouter) Route() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/{url_id}", br.logger.WithLogging(br.handler.DecodeURL))
	router.Post("/", br.logger.WithLogging(br.handler.EncodeURL))

	return router
}
