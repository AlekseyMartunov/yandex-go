package main

import (
	"github.com/AlekseyMartunov/yandex-go.git/internal/app"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebEncodeURL(t *testing.T) {
	a := app.NewApp()

	handler := http.HandlerFunc(a.EncodeURL)
	srv := httptest.NewServer(handler)

	defer srv.Close()

	type wants struct {
		contentType string
		response    string
		statusCode  int
	}

	testCases := []struct {
		name   string
		url    string
		method string
		body   string
		wants  wants
	}{
		{
			name:   "test 1",
			method: http.MethodPost,
			body:   "https://practicum.yandex.ru/",
			url:    "/",
			wants: wants{
				contentType: "text/plain",
				statusCode:  201,
			},
		},

		{
			name:   "test 2",
			method: http.MethodPost,
			body:   "",
			url:    "/",
			wants: wants{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				response:    "Missing body\n",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := resty.New().R()

			req.Method = tc.method

			req.URL = srv.URL + tc.url

			req.SetBody(tc.body)

			resp, err := req.Send()

			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tc.wants.statusCode, resp.StatusCode(),
				"Response code didn't match expected")

			assert.Equal(t, tc.wants.contentType, resp.Header().Get("Content-Type"),
				"Content-Type didn't match expected")

			if tc.wants.response != "" {
				assert.Equal(t, tc.wants.response, string(resp.Body()),
					"Response body didn't match expected")
			}

		})
	}

}

func TestWebDecodeURL(t *testing.T) {
	a := app.NewApp()

	handler := http.HandlerFunc(a.DecodeURL)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	type wants struct {
		contentType string
		response    string
		statusCode  int
		header      string
	}

	testCases := []struct {
		name   string
		url    string
		method string
		wants  wants
	}{
		{
			name:   "test 1",
			method: http.MethodGet,
			url:    "/",
			wants: wants{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				response:    "Empty key\n",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := resty.New().R()

			req.Method = tc.method

			req.URL = srv.URL + tc.url

			resp, err := req.Send()

			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tc.wants.statusCode, resp.StatusCode(),
				"Response code didn't match expected")

			assert.Equal(t, tc.wants.contentType, resp.Header().Get("Content-Type"),
				"Content-Type didn't match expected")

			if tc.wants.response != "" {
				assert.Equal(t, tc.wants.response, string(resp.Body()),
					"Response body didn't match expected")
			}

		})
	}

}
