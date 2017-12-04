package http

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/simonschneider/pefi/middleware"
	"github.com/simonschneider/pefi/services/pefi"
	"net/http"
)

type (
	AccountHandler struct {
		service pefi.AccountService
	}
)

func NewAccountHandler(s pefi.AccountService) *AccountHandler {
	return &AccountHandler{
		service: s,
	}
}

func (h *AccountHandler) Open() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctx := r.Context()
		username := ctx.Value(middleware.Username)
		corrId := ctx.Value(middleware.CorrelationId)
		fmt.Println(username, ",", corrId)
		var acc interface{}
		switch vars["type"] {
		case "external":
			acc, _ = h.service.OpenExternal("accName", "ownerName", "description")
		case "internal":
			acc, _ = h.service.OpenInternal("accName", "ownerName", "description")
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

func (h *AccountHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		account, err := h.service.Get(vars["name"])
		if err != nil {
			fmt.Println("no such account")
			return
		}
		if err := json.NewEncoder(w).Encode(account); err != nil {
			fmt.Println("error encoding")
		}
	}
}

func (h *AccountHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		account, err := h.service.Get(vars["name"])
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

func (h *AccountHandler) Attach(top *mux.Router) {
	router := top.PathPrefix("/api/accounts/").Subrouter()
	router.HandleFunc("/open/{type}", h.Open())
	router.HandleFunc("/{name}", h.Get())
}
