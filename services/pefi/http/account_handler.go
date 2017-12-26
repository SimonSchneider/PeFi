package http

import (
	"context"
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
		ctx := r.Context()
		username := ctx.Value(middleware.Username)
		corrId := ctx.Value(middleware.CorrelationId)
		fmt.Println(username, ",", corrId)
		var acc *pefi.Account
		acc, _ = h.service.Open(context.Background(), "accNameInt", "2ae9052c-8033-477c-a7ae-ae65e6b58879", "description")
		w.Header().Set("test", "this is a test")
		if err := json.NewEncoder(w).Encode(acc); err != nil {
			fmt.Println("error encoding")
		}
	}
}

func (h *AccountHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		account, err := h.service.Get(context.Background(), pefi.ID(vars["name"]))
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
		account, err := h.service.Get(context.Background(), pefi.ID(vars["name"]))
		if err != nil {
			fmt.Println(err)
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

func (h *AccountHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accounts, err := h.service.GetAll(context.Background(), pefi.ID(vars["user"]))
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if accounts == nil {
			msg := "no accounts found"
			fmt.Println(msg)
			http.Error(w, msg, http.StatusNotFound)
			return
		}
		if err := json.NewEncoder(w).Encode(accounts); err != nil {
			fmt.Println("error encoding")
		}
	}
}

func (h *AccountHandler) Attach(top *mux.Router) {
	router := top.PathPrefix("/api/accounts/").Subrouter()
	router.HandleFunc("/open", h.Open()).Name("openAccount").Methods("GET")
	router.HandleFunc("/{name}", h.Get()).Name("getAccount")
	router.HandleFunc("/", h.GetAll()).Queries("user_id", "{user}").Name("getAccounts")
}
