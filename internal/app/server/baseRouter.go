package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ShortURLHandler interface {
	EncodeURL(http.ResponseWriter, *http.Request)
	DecodeURL(http.ResponseWriter, *http.Request)
}

type BaseRouter struct {
	handler ShortURLHandler
}

func NewBaseRouter(h ShortURLHandler) *BaseRouter {
	return &BaseRouter{handler: h}
}

func (br *BaseRouter) Route() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/{url_id}", br.handler.DecodeURL)
	router.Post("/", br.handler.EncodeURL)

	return router
}
