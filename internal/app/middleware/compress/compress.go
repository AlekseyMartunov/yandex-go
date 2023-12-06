// Package compress uses to add function compress in to middleware
package compress

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

var gzipWriter = gzip.NewWriter(&bytes.Buffer{})

// compressWriter mock stub in revenge of the original request
type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// newCompressWriter create new struct
func newCompressWriter(w http.ResponseWriter) *compressWriter {
	gzipWriter.Reset(w)
	return &compressWriter{
		w:  w,
		zw: gzipWriter,
	}
}

// Header uses to override original header
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Write uses to override original body
func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader uses to override original header
func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close uses to close original header
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

// compressReader uses to override original reader
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// newCompressReader return new struct
func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Reade uses to overrider original data
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close uses to close original reader
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// Compress squeezes data
func Compress(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = cr
			defer cr.Close()
		}

		h.ServeHTTP(ow, r)

	})
}
