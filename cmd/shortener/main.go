package main

import (
	"github.com/AlekseyMartunov/yandex-go.git/internal/app"
	"net/http"
)

func main() {
	a := app.NewApp()

	mux := http.NewServeMux()
	mux.HandleFunc("/", a.ShortURLHandler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}

}
