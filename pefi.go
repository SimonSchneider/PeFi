package main

import (
	"context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"os"
	"strconv"
)

func main() {
	router := mux.NewRouter()

	cl := &client{
		address:  ":22400",
		database: &sqlx.DB{},
	}

	endpoints := []endpoint{
		label{cl: cl, extension: "/labels"},
	}

	for _, e := range endpoints {
		router.Handle(e.Extension(), getAll(e)).
			Methods("GET")

		router.Handle(e.Extension(), add(e)).
			Methods("POST")

		router.Handle(e.Extension()+"/{id}",
			idExtraction(get(e))).
			Methods("GET")

		router.Handle(e.Extension()+"/{id}",
			idExtraction(del(e))).
			Methods("DEL")

		router.Handle(e.Extension()+"/{id}",
			idExtraction(mod(e))).
			Methods("PUT")
	}

	http.ListenAndServe(":22400",
		handlers.LoggingHandler(os.Stdout, userAuth(router)))
}

func userAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := strconv.Atoi(r.Header.Get("user"))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
		ctx := context.WithValue(r.Context(), userID, int64(user))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func userFromContext(ctx context.Context) int64 {
	return ctx.Value(userID).(int64)
}

func idExtraction(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}
		ctx := context.WithValue(r.Context(), modelID, int64(id))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
