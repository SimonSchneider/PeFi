package middleware

import (
	"context"
	"net/http"
	"os"
	"time"
)

type (
	ctxKey  int
	Adapter func(http.Handler) http.Handler
)

const (
	Username ctxKey = iota
	CorrelationId
)

func Context(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("username")
		if err != nil {
			return
		}
		if cookie != nil {
			ctx := context.WithValue(r.Context(), Username, cookie.Value)
			ctx = context.WithValue(ctx, CorrelationId, "CorrId")
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func Timer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		w.Header().Set("serverTime", time.Since(start).String())
	})
}

func General(serviceName string) Adapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hostname, _ := os.Hostname()
			w.Header().Set("server", hostname)
			w.Header().Set("service", serviceName)
			next.ServeHTTP(w, r)
		})
	}
}

func Json(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("service", "accountService")
		w.Header().Set("date", time.Now().UTC().Format(time.RFC1123))
		next.ServeHTTP(w, r)
	})
}

func ApplyMiddleware(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}
