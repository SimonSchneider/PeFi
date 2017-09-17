package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type (
	key int

	client struct {
		address  string
		database *sqlx.DB
	}

	endpoint interface {
		Extension() string
		GetNew() interface{}
		GetAll(ctx context.Context) (interface{}, error)
		Add(ctx context.Context, newModel interface{}) error
		Mod(ctx context.Context, id int64, modModel interface{}) error
		Get(ctx context.Context, id int64) (interface{}, error)
		Del(ctx context.Context, id int64) error
	}
)

const (
	userID key = iota
	modelID
)

func getAll(e endpoint) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mods, err := e.GetAll(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(mods); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	})
}

func get(e endpoint) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(modelID).(int64)
		mod, err := e.Get(r.Context(), id)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				w.WriteHeader(http.StatusNotFound)
				break
			default:
				w.WriteHeader(http.StatusBadRequest)
			}
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(mod); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	})
}

func add(e endpoint) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mod := e.GetNew()
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(mod); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		err := e.Add(r.Context(), mod)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				w.WriteHeader(http.StatusNotFound)
				break
			default:
				w.WriteHeader(http.StatusBadRequest)
			}
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func del(e endpoint) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(modelID).(int64)
		err := e.Del(r.Context(), int64(id))
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				w.WriteHeader(http.StatusNotFound)
				break
			default:
				w.WriteHeader(http.StatusBadRequest)
			}
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func mod(e endpoint) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := userFromContext(r.Context())
		mod := e.GetNew()
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(mod); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		err := e.Mod(r.Context(), int64(id), mod)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				w.WriteHeader(http.StatusNotFound)
				break
			default:
				w.WriteHeader(http.StatusBadRequest)
			}
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}
