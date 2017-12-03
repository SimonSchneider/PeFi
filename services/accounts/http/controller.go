package http

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/simonschneider/pefi/middleware"
	"github.com/simonschneider/pefi/services/accounts"
	"net/http"
)

type (
	Handler struct {
		Service accounts.Service
	}
)

func (h *Handler) Open() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctx := r.Context()
		username := ctx.Value(middleware.Username)
		corrId := ctx.Value(middleware.CorrelationId)
		fmt.Println(username, ",", corrId)
		var acc interface{}
		switch vars["type"] {
		case "external":
			acc, _ = h.Service.OpenExternal("accName", "ownerName", "description")
		case "internal":
			acc, _ = h.Service.OpenInternal("accName", "ownerName", "description")
		default:
			msg := "unsupported Type"
			fmt.Println(msg)
			http.Error(w, msg, http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("test", "this is a test")
		if err := json.NewEncoder(w).Encode(acc); err != nil {
			fmt.Println("error encoding")
		}
	}
}

func (h *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		account, err := h.Service.Get(vars["name"])
		if err != nil {
			fmt.Println("no such account")
			return
		}
		if err := json.NewEncoder(w).Encode(account); err != nil {
			fmt.Println("error encoding")
		}
	}
}

func (h *Handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		account, err := h.Service.Get(vars["name"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if account == nil {
			msg := "no such account"
			fmt.Println(msg)
			http.Error(w, msg, http.StatusNotFound)
			return
		}
		if err := json.NewEncoder(w).Encode(account); err != nil {
			fmt.Println("error encoding")
		}
	}
}

func Start(h *Handler) {
	router := mux.NewRouter().PathPrefix("/api/accounts/").Subrouter()
	router.HandleFunc("/open/{type}", h.Open())
	router.HandleFunc("/{name}", h.Get())

	http.ListenAndServe(":8080", middleware.ApplyMiddleware(router,
		middleware.Json,
		middleware.General("accountService"),
		middleware.Timer,
		middleware.Context,
	))
}
