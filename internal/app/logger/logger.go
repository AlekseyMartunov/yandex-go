package logger

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var logger *logrus.Logger = &logrus.Logger{}

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (l *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := l.ResponseWriter.Write(b)
	l.responseData.size = size
	return size, err
}

func (l *loggingResponseWriter) WriteHeader(statusCode int) {
	l.ResponseWriter.WriteHeader(statusCode)
	l.responseData.status = statusCode
}

func Initialize(level logrus.Level) {
	l := logrus.New()
	l.SetLevel(level)
	logger = l
}

func WithLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		responseData := &responseData{}

		lrw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		start := time.Now()
		next.ServeHTTP(&lrw, r)
		logger.Infof("METHOD: %s, URL: %s, TIME %dµs, STATUS: %d, SIZE: %d",
			r.Method, r.RequestURI, time.Since(start)/1000, lrw.responseData.status, lrw.responseData.size)
	}
}
