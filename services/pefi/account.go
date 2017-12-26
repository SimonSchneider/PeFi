package pefi

import (
	"context"
)

type (
	AccountService interface {
		Open(ctx context.Context, name string, ownerID ID, description string) (*Account, error)
		Update(ctx context.Context, id ID, new interface{}) error
		Delete(ctx context.Context, id ID) error
		Get(ctx context.Context, id ID) (*Account, error)
		GetAll(ctx context.Context, userID ID) ([]*Account, error)
		Transfer(ctx context.Context, sender, receiver ID) (string, error)
		Deposit(ctx context.Context, id ID, amount uint64) (string, error)
		Withdraw(ctx context.Context, id ID, amount uint64) (string, error)
	}

	Account struct {
		ID          ID             `json:id`
		Name        string         `json:name`
		OwnerID     ID             `json:owner-id`
		Description string         `json:description`
		Balance     MonetaryAmount `json:amount`
	}
)
