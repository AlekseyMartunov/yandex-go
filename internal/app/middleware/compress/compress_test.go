package compress

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func testHelpHandler(w http.ResponseWriter, r *http.Request) {
	body, error := io.ReadAll(r.Body)
	if error != nil {
		fmt.Println(error)
	}
	w.Write(body)
}

func convert(s string) []byte {
	buf := bytes.Buffer{}
	w := gzip.NewWriter(&buf)

	w.Write([]byte(s))
	w.Close()

	return buf.Bytes()
}

func TestCompress(t *testing.T) {
	srv := httptest.NewServer(Compress(http.HandlerFunc(testHelpHandler)))
	defer srv.Close()

	type wants struct {
		body    []byte
		headers map[string]string
	}

	testCase := []struct {
		name            string
		body            []byte
		acceptEncoding  string
		contentEncoding string
		wants           wants
	}{
		{
			name: "test1",
			body: []byte("some text"),
			wants: wants{
				body: []byte("some text"),
			},
		},

		{
			name:           "test2",
			body:           []byte("some text"),
			acceptEncoding: "gzip",
			wants: wants{
				body: convert("some text"),
			},
		},

		{
			name:            "test3",
			body:            convert("some text"),
			contentEncoding: "gzip",
			wants: wants{
				body: []byte("some text"),
			},
		},

		{
			name:            "test3",
			body:            convert("some text"),
			contentEncoding: "gzip",
			acceptEncoding:  "gzip",
			wants: wants{
				body: convert("some text"),
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {

			client := resty.New()

			resp, err := client.R().
				SetBody(tc.body).
				SetHeader("Accept-Encoding", tc.acceptEncoding).
				SetHeader("Content-Encoding", tc.contentEncoding).
				Post(srv.URL)

			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tc.wants.body, resp.Body())

		})
	}

}
