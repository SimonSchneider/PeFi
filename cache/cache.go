package cache

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"os"
	"pefi/router"
)

type cachingResponseWriter struct {
	status int
	body   []byte
	http.ResponseWriter
}

var conn redis.Conn

func getClient() (redis.Conn, error) {
	if conn != nil {
		return conn, nil
	}
	host := os.Getenv("redis-host")
	port := os.Getenv("redis-port")
	conn, err := redis.Dial("tcp", host+":"+port, redis.DialDatabase(1))
	return conn, err
}

func (w *cachingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *cachingResponseWriter) Write(b []byte) (int, error) {
	w.body = b
	return w.ResponseWriter.Write(b)
}

func HTTPCache(timeout int) router.Adaptor {
	return func(inner http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			rc, err := getClient()
			if err == nil {
				val, err := redis.String(rc.Do("GET", r.RequestURI))
				if err == nil {
					w.Write([]byte(val))
					return
				}
			}
			cachingRW := &cachingResponseWriter{
				ResponseWriter: w,
			}
			inner.ServeHTTP(cachingRW, r)
			if cachingRW.status != http.StatusOK {
				fmt.Println("not caching", cachingRW.status, http.StatusOK)
				return
			}
			_, err = rc.Do("SET", r.RequestURI, cachingRW.body)
			if err != nil {
				return
			}
			if timeout != 0 {
				rc.Do("EXPIRE", r.RequestURI, timeout)
			}
		}
	}
}

func HTTPWipe(key string) router.Adaptor {
	return func(inner http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			rc, err := getClient()
			if err == nil {
				rc.Do("DEL", key)
				rc.Do("DEL", r.RequestURI)
			}
			inner.ServeHTTP(w, r)
		}
	}
}
