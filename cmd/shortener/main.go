package main

import (
	"net/http"

	"github.com/AlekseyMartunov/yandex-go.git/internal/app/config"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/encoder"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/handlers"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/logger"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/server"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/storage"
)

func main() {
	c := config.NewConfig()
	c.GetConfig()

	s := storage.NewStorage()
	e := encoder.NewEncoder(s)

	h := handlers.NewShortURLHandler(e, c)
	l := logger.NewLogger("info")
	r := server.NewBaseRouter(h, l)

	err := http.ListenAndServe(c.GetAddress(), r.Route())
	if err != nil {
		panic(err)
	}

}
