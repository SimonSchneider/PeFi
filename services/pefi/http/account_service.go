package http

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/simonschneider/pefi/services/pefi"
	"net/http"
)

type (
	AccountService struct {
		r *mux.Router
		c *http.Client
	}
)

func NewAccountService(router *mux.Router) *AccountService {
	client := &http.Client{}
	return &AccountService{
		r: router,
		c: client,
	}
}

func (a *AccountService) Open(ctx context.Context, name, owner, description string) (*pefi.Account, error) {
	url, err := a.r.Get("openAccount").URL("type", "internal")
	if err != nil {
		return nil, err
	}
	req, err := getRequest(url)
	if err != nil {
		return nil, err
	}
	resp, err := a.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	acc := &pefi.Account{}
	err = json.NewDecoder(resp.Body).Decode(acc)
	if err != nil {
		return nil, err
	}
	return acc, err
}

func (a *AccountService) Update(ctx context.Context, name string, new interface{}) error {
	return nil
}

func (a *AccountService) Delete(ctx context.Context, name string) error {
	return nil
}
func (a *AccountService) Get(ctx context.Context, name string) (*pefi.Account, error) {
	url, err := a.r.Get("getAccount").URL("name", name)
	if err != nil {
		return nil, err
	}
	req, err := getRequest(url)
	if err != nil {
		return nil, err
	}
	resp, err := a.c.Do(req)
	if err != nil {
		return nil, err
	}
	var acc *pefi.Account
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&acc); err != nil {
		return nil, err
	}
	return acc, nil
}
func (a *AccountService) Transfer(ctx context.Context, sender, receiver string) (string, error) {
	return "done", nil
}
func (a *AccountService) Deposit(ctx context.Context, name string, amount uint64) (string, error) {
	return "done", nil
}
func (a *AccountService) Withdraw(ctx context.Context, name string, amount uint64) (string, error) {
	return "done", nil
}
