package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ShortURLHandler interface {
	EncodeURL(w http.ResponseWriter, r *http.Request)
	DecodeURL(w http.ResponseWriter, r *http.Request)
}

type APIHandler interface {
	EncodeAPI(w http.ResponseWriter, r *http.Request)
}

type BaseRouter struct {
	handler    ShortURLHandler
	apiHandler APIHandler
	middleware []func(handler http.Handler) http.Handler
}

func NewBaseRouter(h ShortURLHandler, ah APIHandler, m ...func(handler http.Handler) http.Handler) *BaseRouter {
	return &BaseRouter{
		handler:    h,
		apiHandler: ah,
		middleware: m,
	}
}

func (br *BaseRouter) Route() *chi.Mux {
	router := chi.NewRouter()
	router.Use(br.middleware...)
	router.Get("/{url_id}", br.handler.DecodeURL)
	router.Post("/", br.handler.EncodeURL)
	router.Post("/api/shorten", br.apiHandler.EncodeAPI)

	return router
}
