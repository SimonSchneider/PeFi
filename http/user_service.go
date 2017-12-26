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
	UserService struct {
		r *mux.Router
		c *http.Client
	}
)

func NewUserService(router *mux.Router) *UserService {
	client := &http.Client{}
	return &UserService{
		r: router,
		c: client,
	}
}

func (a *UserService) Create(ctx context.Context, name string) (*pefi.User, error) {
	url, err := a.r.Get("createUser").URL("type", "internal")
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
	user := &pefi.User{}
	err = json.NewDecoder(resp.Body).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (a *UserService) Update(ctx context.Context, id pefi.ID, new interface{}) error {
	return nil
}

func (a *UserService) Delete(ctx context.Context, id pefi.ID) error {
	return nil
}
func (a *UserService) Get(ctx context.Context, id pefi.ID) (*pefi.User, error) {
	url, err := a.r.Get("getUser").URL("id", string(id))
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
	var user *pefi.User
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, errors.New("No such account")
	}
	return user, nil
}
