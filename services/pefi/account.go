package pefi

import (
	"context"
)

type (
	AccountService interface {
		Open(ctx context.Context, name, owner, description string) (*Account, error)
		Update(ctx context.Context, name string, new interface{}) error
		Delete(ctx context.Context, name string) error
		Get(ctx context.Context, name string) (*Account, error)
		Transfer(ctx context.Context, sender, receiver string) (string, error)
		Deposit(ctx context.Context, name string, amount uint64) (string, error)
		Withdraw(ctx context.Context, name string, amount uint64) (string, error)
	}

	MonetaryAmount struct {
		Amount   int64  `json:amount`
		Currency string `json:currency`
	}

	Account struct {
		Name        string         `json:name`
		OwnerId     string         `json:owner-id`
		Description string         `json:description`
		Amount      MonetaryAmount `json:amount`
	}
)
