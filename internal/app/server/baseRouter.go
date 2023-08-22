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

type APIHandler interface {
	EncodeAPI(w http.ResponseWriter, r *http.Request)
}

type BaseRouter struct {
	handler    ShortURLHandler
	logger     Logger
	apiHandler APIHandler
}

func NewBaseRouter(h ShortURLHandler, ah APIHandler, l Logger) *BaseRouter {
	return &BaseRouter{
		handler:    h,
		apiHandler: ah,
		logger:     l,
	}
}

func (br *BaseRouter) Route() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/{url_id}", br.logger.WithLogging(br.handler.DecodeURL))
	router.Post("/", br.logger.WithLogging(br.handler.EncodeURL))
	router.Post("/api/shorten", br.logger.WithLogging(br.apiHandler.EncodeAPI))

	return router
}
