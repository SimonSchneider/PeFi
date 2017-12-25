package redis

import (
	"context"
	"fmt"
	"github.com/simonschneider/pefi/services/pefi"
)

type (
	AccountService struct {
		cache   map[string]*pefi.Account
		service pefi.AccountService
	}
)

func NewAccountService(next pefi.AccountService) *AccountService {
	return &AccountService{
		cache:   make(map[string]*pefi.Account),
		service: next,
	}
}

func (s AccountService) Open(ctx context.Context, name, owner, description string) (*pefi.Account, error) {
	a, err := s.service.Open(ctx, name, owner, description)
	if err != nil {
		return a, err
	}
	s.cache[a.Name] = a
	fmt.Println("added to cache" + a.Name)
	return a, err
}

func (s AccountService) Update(ctx context.Context, name string, new interface{}) error {
	return s.service.Update(ctx, name, new)
}

func (s AccountService) Delete(ctx context.Context, name string) error {
	return s.service.Delete(ctx, name)
}

func (s AccountService) Get(ctx context.Context, name string) (*pefi.Account, error) {
	if acc := s.cache[name]; acc != nil {
		return acc, nil
	}
	return s.service.Get(ctx, name)
}

func (s AccountService) Transfer(ctx context.Context, sender, receiver string) (string, error) {
	return "transfer id", nil
}

func (s AccountService) Deposit(ctx context.Context, name string, amount uint64) (string, error) {
	return "deposit id", nil
}

func (s AccountService) Withdraw(ctx context.Context, name string, amount uint64) (string, error) {
	return "withdraw id", nil
}
