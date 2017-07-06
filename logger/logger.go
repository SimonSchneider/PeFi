package logger

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"pefi/router"
	"time"
)

type loggingResponseWriter struct {
	status int
	http.ResponseWriter
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func HTTPLogger(info string) router.HTTPDecorator {
	return func(inner http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			loggingRW := &loggingResponseWriter{
				ResponseWriter: w,
			}
			start := time.Now()
			inner.ServeHTTP(loggingRW, r)
			requestLogger := log.WithFields(
				log.Fields{
					"method": r.Method,
					"url":    r.RequestURI,
					"time":   time.Since(start),
					"status": loggingRW.status,
				})
			requestLogger.Info(info)
		}
	}
}
