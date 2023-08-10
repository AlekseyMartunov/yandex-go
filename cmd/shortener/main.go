package main

import (
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/config"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/handlers"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/server"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/storage"
	"net/http"
)

func main() {
	s := storage.NewStorage()
	c := config.NewConfig()
	c.GetConfig()
	h := handlers.NewShortUrlHandler(s, c)
	r := server.NewBaseRouter(h)

	err := http.ListenAndServe(c.GetAddress(), r.Route())
	if err != nil {
		panic(err)
	}

}
