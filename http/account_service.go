package http

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/simonschneider/pefi"
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

func (a *AccountService) Open(ctx context.Context, name string, owner pefi.ID, description string) (*pefi.Account, error) {
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

func (a *AccountService) Update(ctx context.Context, id pefi.ID, new interface{}) error {
	return nil
}

func (a *AccountService) Delete(ctx context.Context, id pefi.ID) error {
	return nil
}
func (a *AccountService) Get(ctx context.Context, id pefi.ID) (*pefi.Account, error) {
	url, err := a.r.Get("getAccount").URL("name", string(id))
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
		return nil, errors.New("No such account")
	}
	return acc, nil
}

func (a *AccountService) GetAll(ctx context.Context, userID pefi.ID) ([]*pefi.Account, error) {
	url, err := a.r.Get("getAccounts").URL("user", string(userID))
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
	var accs []*pefi.Account
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&accs); err != nil {
		return nil, errors.New("No such account")
	}
	return accs, nil
}

func (a *AccountService) Transfer(ctx context.Context, sender, receiver pefi.ID) (string, error) {
	return "done", nil
}
func (a *AccountService) Deposit(ctx context.Context, name pefi.ID, amount uint64) (string, error) {
	return "done", nil
}
func (a *AccountService) Withdraw(ctx context.Context, name pefi.ID, amount uint64) (string, error) {
	return "done", nil
}
