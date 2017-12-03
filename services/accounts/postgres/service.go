package postgres

import (
	"errors"
	"github.com/simonschneider/pefi/services/accounts"
)

type (
	Service struct {
	}
)

func (s Service) OpenExternal(name, owner, description string) (*accounts.ExternalAccount, error) {
	return &accounts.ExternalAccount{
		Name:        name,
		OwnerId:     owner,
		Description: description,
	}, nil
}

func (s Service) OpenInternal(name, owner, description string) (*accounts.InternalAccount, error) {
	return &accounts.InternalAccount{
		ExternalAccount: accounts.ExternalAccount{
			Name:        name,
			OwnerId:     owner,
			Description: description,
		},
		Amount: accounts.MonetaryAmount{0, "SEK"},
	}, nil
}

func (s Service) Update(name string, new interface{}) error {
	return nil
}

func (s Service) Delete(name string) error {
	return nil
}

func (s Service) Get(name string) (interface{}, error) {
	return nil, errors.New("no such account")
}

func (s Service) Transfer(sender, receiver string) (string, error) {
	return "transfer id", nil
}

func (s Service) Deposit(name string, amount uint64) (string, error) {
	return "deposit id", nil
}

func (s Service) Withdraw(name string, amount uint64) (string, error) {
	return "withdraw id", nil
}

func SaveExternal(a accounts.ExternalAccount) error {
	//_, err := s.db.Exec("INSERT INTO external_accounts(name, owner-id, description) VALUES($1, $2, $3)", a.Name, a.OwnerId, a.Description)
	return nil
}

func SaveInternal(a accounts.InternalAccount) error {
	//_, err := s.db.Exec("INSERT INTO internal_accounts(name, owner-id, description, amount) VALUES($1, $2, $3, $4)", a.Name, a.OwnerId, a.Description, a.Amount)
	return nil
}
