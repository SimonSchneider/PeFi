package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/simonschneider/pefi"
	"github.com/simonschneider/pefi/middleware"
	"net/http"
)

type (
	UserHandler struct {
		service pefi.UserService
	}
)

func NewUserHandler(s pefi.UserService) *UserHandler {
	return &UserHandler{
		service: s,
	}
}

func (h *UserHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		username := ctx.Value(middleware.Username)
		corrId := ctx.Value(middleware.CorrelationId)
		fmt.Println(username, ",", corrId)
		var acc *pefi.User
		acc, _ = h.service.Create(context.Background(), "user")
		w.Header().Set("test", "this is a test")
		if err := json.NewEncoder(w).Encode(acc); err != nil {
			fmt.Println("error encoding")
		}
	}
}

func (h *UserHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		user, err := h.service.Get(context.Background(), pefi.ID(vars["id"]))
		if err != nil {
			fmt.Println("no such user")
			return
		}
		if err := json.NewEncoder(w).Encode(user); err != nil {
			fmt.Println("error encoding")
		}
	}
}

func (h *UserHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		user, err := h.service.Get(context.Background(), pefi.ID(vars["id"]))
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if user == nil {
			msg := "no such user"
			fmt.Println(msg)
			http.Error(w, msg, http.StatusNotFound)
			return
		}
		if err := json.NewEncoder(w).Encode(user); err != nil {
			fmt.Println("error encoding")
		}
	}
}

func (h *UserHandler) Attach(top *mux.Router) {
	router := top.PathPrefix("/api/users/").Subrouter()
	router.HandleFunc("/create", h.Create()).Name("createUser").Methods("GET")
	router.HandleFunc("/{id}", h.Get()).Name("getUser")
}
