package api

import (
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"

	config2 "github.com/AlekseyMartunov/yandex-go.git/internal/app/config"
	encoder2 "github.com/AlekseyMartunov/yandex-go.git/internal/app/encoder"
	storage2 "github.com/AlekseyMartunov/yandex-go.git/internal/app/storage"
)

func TestApiHandlers(t *testing.T) {
	s := storage2.NewStorage()
	c := config2.NewConfig()

	e := encoder2.NewEncoder(s)
	h := NewApiHandlers(e, c)

	handler := http.HandlerFunc(h.EncodeAPI)
	srv := httptest.NewServer(handler)

	defer srv.Close()

	type wants struct {
		statusCode  int
		body        string
		contentType string
	}

	testCases := []struct {
		name        string
		url         string
		body        string
		contentType string
		wants       wants
	}{
		{
			name:        "test1",
			url:         "/api/shorten",
			body:        `{"url": "123AAA"}`,
			contentType: "application/json",
			wants: wants{
				statusCode:  http.StatusCreated,
				contentType: "application/json",
			},
		},

		{
			name:        "test2",
			url:         "/api/shorten",
			body:        `{"url: "123AAA"}`,
			contentType: "application/json",
			wants: wants{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "Request body read error\n",
			},
		},

		{
			name:        "test3",
			url:         "/api/shorten",
			body:        `{"url": "123AAA"}`,
			contentType: "text/plain",
			wants: wants{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "Content-type should be 'application/json'\n",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			req := resty.New().R()
			req.Method = http.MethodPost
			req.Body = tc.body
			req.URL = srv.URL + tc.url
			req.Header.Set("Content-Type", tc.contentType)

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tc.wants.statusCode, resp.StatusCode(),
				"Response code didn't match expected")

			assert.Equal(t, tc.wants.contentType, resp.Header().Get("Content-Type"),
				"Content-Type  didn't match expected")

			if tc.wants.body != "" {
				assert.Equal(t, tc.wants.body, string(resp.Body()),
					"Response body didn't match expected")
			}

		})
	}

}
