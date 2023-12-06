package handlers

import (
	"context"
	storage2 "github.com/AlekseyMartunov/yandex-go.git/internal/app/model/url/simpleurl"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	config2 "github.com/AlekseyMartunov/yandex-go.git/internal/app/config"
	encoder2 "github.com/AlekseyMartunov/yandex-go.git/internal/app/encoder"
)

func TestShortUrlHandler_EncodeURL(t *testing.T) {
	s, _ := storage2.NewMapStorage()
	ctx := context.Background()

	c := config2.NewConfig()
	//c.GetConfig()
	e := encoder2.NewEncoder(s)
	h := NewShortURLHandler(e, c, ctx)

	handler := http.HandlerFunc(h.EncodeURL)
	srv := httptest.NewServer(handler)

	defer srv.Close()

	type wants struct {
		statusCode  int
		body        string
		contentType string
	}

	testCases := []struct {
		name  string
		url   string
		body  string
		wants wants
	}{
		{
			name: "test1",
			url:  "/",
			body: "someURL",
			wants: wants{
				statusCode:  http.StatusCreated,
				body:        "should be no empty body",
				contentType: "text/plain",
			},
		},
		{
			name: "test2",
			url:  "/",
			wants: wants{
				statusCode:  http.StatusBadRequest,
				body:        "Missing body\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			req := resty.New().R()
			req.Method = http.MethodPost
			req.Body = tc.body
			req.URL = srv.URL + tc.url

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tc.wants.statusCode, resp.StatusCode(),
				"Response code didn't match expected")

			assert.Equal(t, tc.wants.contentType, resp.Header().Get("Content-Type"),
				"Response content-type didn't match expected")

			if tc.wants.body == "should be no empty body" {
				assert.NotEmpty(t, string(resp.Body()),
					"Response body didn't match expected")
			}
		})
	}

}

func TestShortUrlHandler_DecodeURL(t *testing.T) {
	s, _ := storage2.NewMapStorage()
	ctx := context.Background()

	c := config2.NewConfig()
	//c.GetConfig()
	e := encoder2.NewEncoder(s)
	h := NewShortURLHandler(e, c, ctx)

	handler := http.HandlerFunc(h.DecodeURL)
	srv := httptest.NewServer(handler)

	defer srv.Close()

	type wants struct {
		statusCode  int
		body        string
		contentType string
	}

	testCases := []struct {
		name  string
		url   string
		wants wants
	}{
		{
			name: "test1",
			url:  "/",
			wants: wants{
				statusCode:  http.StatusBadRequest,
				body:        "Empty key\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "test2",
			url:  "/someString",
			wants: wants{
				statusCode:  http.StatusBadRequest,
				body:        "Empty key\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			req := resty.New().R()
			req.Method = http.MethodGet
			req.URL = srv.URL + tc.url

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tc.wants.statusCode, resp.StatusCode(),
				"Response code didn't match expected")

			assert.Equal(t, tc.wants.contentType, resp.Header().Get("Content-Type"),
				"Response content-type didn't match expected")

			assert.Equal(t, tc.wants.body, string(resp.Body()),
				"Response body didn't match expected")
		})
	}

}
