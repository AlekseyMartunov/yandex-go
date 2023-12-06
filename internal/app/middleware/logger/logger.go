// Package logger provider logging in application
package logger

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger simple logging struct
type Logger struct {
	defaultLogger *logrus.Logger
}

// NewLogger create new logger instance
func NewLogger(level string) *Logger {
	l := logrus.New()
	lvl, err := logrus.ParseLevel(level)

	if err != nil {
		panic(err)
	}

	l.SetLevel(lvl)
	return &Logger{defaultLogger: l}
}

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

// Write provides Write method to mock
func (l *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := l.ResponseWriter.Write(b)
	l.responseData.size = size
	return size, err
}

// WriteHeader provides WriteHeader method to mock
func (l *loggingResponseWriter) WriteHeader(statusCode int) {
	l.ResponseWriter.WriteHeader(statusCode)
	l.responseData.status = statusCode
}

// WithLogging function to logging requests in middleware
func (l *Logger) WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		responseData := &responseData{}

		lrw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		start := time.Now()
		next.ServeHTTP(&lrw, r)
		l.defaultLogger.Infof("METHOD: %s, URL: %s, TIME %dµs, STATUS: %d, SIZE: %d",
			r.Method, r.RequestURI, time.Since(start)/1000, lrw.responseData.status, lrw.responseData.size)
	})
}
