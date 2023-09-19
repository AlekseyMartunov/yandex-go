package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ShortURLHandler interface {
	EncodeURL(w http.ResponseWriter, r *http.Request)
	DecodeURL(w http.ResponseWriter, r *http.Request)
	EncodeAPI(w http.ResponseWriter, r *http.Request)
	DataBaseStatus(w http.ResponseWriter, r *http.Request)
	BatchURL(w http.ResponseWriter, r *http.Request)
	GetAllURL(w http.ResponseWriter, r *http.Request)
}

type BaseRouter struct {
	handler    ShortURLHandler
	middleware []func(handler http.Handler) http.Handler
}

func NewBaseRouter(h ShortURLHandler, m ...func(handler http.Handler) http.Handler) *BaseRouter {
	return &BaseRouter{
		handler:    h,
		middleware: m,
	}
}

func (br *BaseRouter) Route() *chi.Mux {
	router := chi.NewRouter()
	router.Use(br.middleware...)

	router.Get("/{url_id}", br.handler.DecodeURL)
	router.Get("/ping", br.handler.DataBaseStatus)
	router.Get("/api/user/urls", br.handler.GetAllURL)

	router.Post("/", br.handler.EncodeURL)
	router.Post("/api/shorten", br.handler.EncodeAPI)
	router.Post("/api/shorten/batch", br.handler.BatchURL)

	return router
}
