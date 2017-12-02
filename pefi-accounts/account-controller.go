package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/simonschneider/pefi/pefi-accounts/middleware"
	"net/http"
)

type (
	AccountController struct {
		accountService *AccountService
	}
)

func NewAccountController(accountService *AccountService) *AccountController {
	return &AccountController{
		accountService: accountService,
	}
}

func (a AccountController) Open() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctx := r.Context()
		username := ctx.Value(middleware.Username)
		corrId := ctx.Value(middleware.CorrelationId)
		fmt.Println(username, ",", corrId)
		var acc Account
		switch vars["type"] {
		case "external":
			acc, _ = a.accountService.OpenExternal("accName", "ownerName", "description")
		case "internal":
			acc, _ = a.accountService.OpenInternal("accName", "ownerName", "description")
		default:
			fmt.Println("unsupported Type")
			return
		}
		w.Header().Set("test", "this is a test")
		if err := json.NewEncoder(w).Encode(acc); err != nil {
			fmt.Println("error encoding")
		}
	}
}
