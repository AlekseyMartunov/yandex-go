// Package router store endpoints
package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ShortURLHandler interface show all endpoints in handlers
type ShortURLHandler interface {
	EncodeURL(w http.ResponseWriter, r *http.Request)
	DecodeURL(w http.ResponseWriter, r *http.Request)
	EncodeAPI(w http.ResponseWriter, r *http.Request)
	DataBaseStatus(w http.ResponseWriter, r *http.Request)
	BatchURL(w http.ResponseWriter, r *http.Request)
	GetAllURL(w http.ResponseWriter, r *http.Request)
	DeleteURL(w http.ResponseWriter, r *http.Request)
	Stats(w http.ResponseWriter, r *http.Request)
}

// BaseRouter struct applications router
type BaseRouter struct {
	handler    ShortURLHandler
	middleware []func(handler http.Handler) http.Handler
}

// NewBaseRouter create new instance
func NewBaseRouter(h ShortURLHandler, m ...func(handler http.Handler) http.Handler) *BaseRouter {
	return &BaseRouter{
		handler:    h,
		middleware: m,
	}
}

// Route return chi router
func (br *BaseRouter) Route() *chi.Mux {
	router := chi.NewRouter()
	router.Use(br.middleware...)

	router.Get("/{url_id}", br.handler.DecodeURL)
	router.Get("/ping", br.handler.DataBaseStatus)
	router.Get("/api/user/urls", br.handler.GetAllURL)
	router.Get("/api/internal/stats", br.handler.Stats)

	router.Post("/", br.handler.EncodeURL)
	router.Post("/api/shorten", br.handler.EncodeAPI)
	router.Post("/api/shorten/batch", br.handler.BatchURL)

	router.Delete("/api/user/urls", br.handler.DeleteURL)

	return router
}
