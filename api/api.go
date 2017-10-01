package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/simonschneider/pefi/api/models"
	"net/http"
	"os"
	"strconv"
)

var (
	db *sqlx.DB
)

type (
	endpoint interface {
		URL() string
		GetNewInstance() interface{}
		GetAll(user int64) (interface{}, error)
		Add(user int64, IDENT interface{}) error
		Get(user, id int64) (interface{}, error)
		Del(user, id int64) error
		Mod(user, id int64, IDENT interface{}) error
	}
)

func main() {
	router := mux.NewRouter()

	endpoints := []models.Endpoint{
		&models.Category{},
		&models.ExternalAccount{},
		&models.InternalAccount{},
		&models.Label{},
		&models.Transaction{},
	}

	for _, e := range endpoints {
		router.
			Handle(e.URL(), getAll(e)).
			Methods("GET")

		router.
			Handle(e.URL(), add(e)).
			Methods("POST")

		router.
			Handle(e.URL()+"/{id}", get(e)).
			Methods("GET")

		router.
			Handle(e.URL()+"/{id}", del(e)).
			Methods("DEL")

		router.
			Handle(e.URL()+"/{id}", mod(e)).
			Methods("PUT")
	}

	router.
		Handle("/testing", Tmp()).
		Methods("GET")

	router.
		Handle("/testing/{id}", TmpG()).
		Methods("GET")

	dbHost := os.Getenv("postgres-host")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbinfo := "host=" + dbHost + " user=postgres database=pefi sslmode=disable"

	tmp, err := sqlx.Connect("postgres", dbinfo)
	if err != nil {
		fmt.Println("error1")
		return
	}
	if err = tmp.Ping(); err != nil {
		fmt.Println("error2")
		return
	}
	db = tmp
	models.InitDB(db)

	Init(dbinfo)

	fmt.Println("starting")

	http.ListenAndServe(":22400", handlers.LoggingHandler(os.Stdout, router))
}

func getAll(e models.Endpoint) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := strconv.Atoi(r.Header.Get("user"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		models, err := e.GetAll(int64(user))
		fmt.Println(err)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err = json.NewEncoder(w).Encode(models); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func get(e models.Endpoint) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		user, err := strconv.Atoi(r.Header.Get("user"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		mod, err := e.Get(int64(user), int64(id))
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				w.WriteHeader(http.StatusNotFound)
			default:
				w.WriteHeader(http.StatusBadRequest)
			}
			return
		}
		if err = json.NewEncoder(w).Encode(mod); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
	})
}

func add(e models.Endpoint) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := strconv.Atoi(r.Header.Get("user"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(e); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = e.Add(int64(user))
		fmt.Println(err)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				w.WriteHeader(http.StatusNotFound)
			default:
				w.WriteHeader(http.StatusBadRequest)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func del(e models.Endpoint) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		user, err := strconv.Atoi(r.Header.Get("user"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = e.Del(int64(user), int64(id))
		fmt.Println(err)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				w.WriteHeader(http.StatusNotFound)
			default:
				w.WriteHeader(http.StatusBadRequest)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func mod(e models.Endpoint) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		user, err := strconv.Atoi(r.Header.Get("user"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(e); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = e.Mod(int64(user), int64(id))
		fmt.Println(err)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				w.WriteHeader(http.StatusNotFound)
			default:
				w.WriteHeader(http.StatusBadRequest)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func mwHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}
