package main

import (
	"fmt"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/config"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/server"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWebDecodeURL(t *testing.T) {

	s := storage.NewStorage()
	c := config.NewConfig()
	r := server.NewBaseRouter(s, c)
	router := r.Route()

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
			wants: wants{
				statusCode:  http.StatusMethodNotAllowed,
				body:        "",
				contentType: "",
			},
		},
		{
			name: "test1",
			url:  "/jgkflf",
			wants: wants{
				statusCode:  http.StatusBadRequest,
				body:        "Empty key\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			req := httptest.NewRequest("GET", tc.url, nil)
			respRec := httptest.NewRecorder()
			router.ServeHTTP(respRec, req)

			assert.Equal(t, tc.wants.statusCode, respRec.Code,
				"Response code didn't match expected")

			assert.Equal(t, tc.wants.body, fmt.Sprintf("%v", respRec.Body),
				"Response body didn't match expected")

			assert.Equal(t, tc.wants.contentType, respRec.Header().Get("Content-Type"))

		})
	}

}

func TestWebEncodeURL(t *testing.T) {

	s := storage.NewStorage()
	c := config.NewConfig()
	r := server.NewBaseRouter(s, c)
	router := r.Route()

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

			req := httptest.NewRequest("POST", tc.url, strings.NewReader(tc.body))
			respRec := httptest.NewRecorder()
			router.ServeHTTP(respRec, req)

			assert.Equal(t, tc.wants.statusCode, respRec.Code,
				"Response code didn't match expected")

			assert.Equal(t, tc.wants.contentType, respRec.Header().Get("Content-Type"))

			if tc.wants.body == "should be no empty body" {
				assert.NotEmpty(t, respRec.Body,
					"Response body should be not empty")
			} else {
				assert.Equal(t, tc.wants.body, fmt.Sprintf("%v", respRec.Body),
					"Response body didn't match expected")
			}

		})
	}

}
