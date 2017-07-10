package cache

import (
	"github.com/garyburd/redigo/redis"
	"net/http"
)

type cachingResponseWriter struct {
	status int
	body   []byte
	http.ResponseWriter
}

func (w *cachingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *cachingResponseWriter) Write(b []byte) (int, error) {
	w.body = b
	return w.ResponseWriter.Write(b)
}

// HTTPCache the response of the request if it returns StatusOK,
// the request URL is the key that will be matched
func HTTPCache(c *Client, timeout int) func(http.HandlerFunc) http.HandlerFunc {
	return func(inner http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			val, err := redis.String(c.Do("GET", r.RequestURI))
			if err == nil {
				w.Write([]byte(val))
				return
			}
			cachingRW := &cachingResponseWriter{
				ResponseWriter: w,
			}
			inner.ServeHTTP(cachingRW, r)
			if cachingRW.status != http.StatusOK {
				return
			}
			_, err = c.Do("SET", r.RequestURI, cachingRW.body)
			if err != nil {
				return
			}
			if timeout != 0 {
				c.Do("EXPIRE", r.RequestURI, timeout)
			}
		}
	}
}

// HTTPWipeCache wipes the keys from the cache including the Requested
// URL
func HTTPWipeCache(c *Client, key ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(inner http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			for _, k := range key {
				c.Do("DEL", k)
			}
			c.Do("DEL", r.RequestURI)
			inner.ServeHTTP(w, r)
		}
	}
}
