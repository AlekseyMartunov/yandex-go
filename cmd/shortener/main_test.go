package main

import (
	"fmt"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShortURLHandler(t *testing.T) {
	a := app.NewApp()

	type wants struct {
		contentType string
		response    string
		statusCode  int
	}

	tests := []struct {
		name   string
		url    string
		method string
		body   string
		wants  wants
	}{
		{
			name:   "test 1",
			url:    "/",
			method: http.MethodPost,
			body:   "https://practicum.yandex.ru/",
			wants: wants{
				contentType: "text/plain",
				statusCode:  201,
			},
		},

		{
			name:   "test 2",
			url:    "/Ejfkdsh",
			method: http.MethodPost,
			body:   "https://practicum.yandex.ru/",
			wants: wants{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				response:    "You should send a request to '/'\n",
			},
		},

		{
			name:   "test 3",
			url:    "/Ejfkdsh",
			method: http.MethodGet,
			wants: wants{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				response:    "Empty key\n",
			},
		},

		{
			name:   "test 4",
			url:    "/",
			method: http.MethodPost,
			wants: wants{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				response:    "Missing body\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			a.ShortURLHandler(w, r)

			res := w.Result()

			assert.Equal(t, tt.wants.statusCode, res.StatusCode)
			assert.Equal(t, tt.wants.contentType, res.Header.Get("Content-Type"))

			defer res.Body.Close()

			if tt.wants.response != "" {
				resBody, err := io.ReadAll(res.Body)
				require.NoError(t, err)

				assert.Equal(t, tt.wants.response, string(resBody))
			}

			fmt.Println()

		})

	}

}
