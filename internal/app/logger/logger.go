package logger

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	defaultLogger *logrus.Logger
}

func NewLogger() *Logger {
	return &Logger{}
}

var logger Logger = Logger{defaultLogger: &logrus.Logger{}}

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

func init() {
	l := logrus.New()
	l.SetLevel(logrus.InfoLevel)
	logger.defaultLogger = l
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
		logger.defaultLogger.Infof("METHOD: %s, URL: %s, TIME %dµs, STATUS: %d, SIZE: %d",
			r.Method, r.RequestURI, time.Since(start)/1000, lrw.responseData.status, lrw.responseData.size)
	}
}
