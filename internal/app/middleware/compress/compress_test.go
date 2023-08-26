package compress

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testHelpHandler(w http.ResponseWriter, r *http.Request) {
	body, error := ioutil.ReadAll(r.Body)
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

func deconvert(b []byte) string {
	buf := bytes.NewBuffer(b)
	r, err := gzip.NewReader(buf)

	defer r.Close()

	if err != nil {
		fmt.Println(err)
	}

	res, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Println(err)
	}
	return string(res)
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

			fmt.Println("body is", string(resp.Body()))

			assert.Equal(t, tc.wants.body, resp.Body())

		})
	}

}
