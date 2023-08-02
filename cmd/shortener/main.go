package main

import (
	"github.com/AlekseyMartunov/yandex-go.git/internal/app"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	a := app.NewApp()

	r := chi.NewRouter()
	r.Get("/{url_id}", a.DecodeURL)
	r.Post("/", a.EncodeURL)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}

}
