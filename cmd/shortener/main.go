package main

import (
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/config"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/handlers"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type App struct {
	server *handlers.Server
}

func main() {

	s := storage.NewStorage()
	cfg := config.Config()

	server := handlers.NewServer(s, cfg)

	a := &App{
		server: server,
	}

	r := chi.NewRouter()
	r.Get("/{url_id}", a.server.DecodeURL)
	r.Post("/", a.server.EncodeURL)

	err := http.ListenAndServe(a.server.Cfg.Host, r)
	if err != nil {
		panic(err)
	}

}
