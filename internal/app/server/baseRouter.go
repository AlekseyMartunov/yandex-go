package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/AlekseyMartunov/yandex-go.git/internal/app/logger"
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
	router.Get("/{url_id}", logger.WithLogging(br.handler.DecodeURL))
	router.Post("/", logger.WithLogging(br.handler.EncodeURL))

	return router
}
