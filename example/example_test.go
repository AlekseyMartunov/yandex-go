package example

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/AlekseyMartunov/yandex-go.git/internal/app/config"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/encoder"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/handlers"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/url/simpleurl"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/router"
)

func Example() {
	storage, err := simpleurl.NewMapStorage()
	if err != nil {
		log.Fatalln(err)
	}
	cfg := config.NewConfig()
	encoder := encoder.NewEncoder(storage)
	handler := handlers.NewShortURLHandler(encoder, cfg, context.Background())
	router := router.NewBaseRouter(handler)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://practicum.yandex.ru/ "))
	w := httptest.NewRecorder()

	router.Route().ServeHTTP(w, req)

	shortedURL := w.Body.String()

	req = httptest.NewRequest(http.MethodGet, shortedURL, nil)
	w = httptest.NewRecorder()

	router.Route().ServeHTTP(w, req)

	fmt.Println(w.Header().Get("Location"))

	//Output:
	//https://practicum.yandex.ru/

}
