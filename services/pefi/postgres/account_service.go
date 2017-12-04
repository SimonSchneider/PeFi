package postgres

import (
	"errors"
	"github.com/simonschneider/pefi/services/pefi"
)

type (
	AccountService struct {
	}
)

func NewAccountService() *AccountService {
	return &AccountService{}
}

func (s AccountService) OpenExternal(name, owner, description string) (*pefi.ExternalAccount, error) {
	return &pefi.ExternalAccount{
		Name:        name,
		OwnerId:     owner,
		Description: description,
	}, nil
}

func (s AccountService) OpenInternal(name, owner, description string) (*pefi.InternalAccount, error) {
	return &pefi.InternalAccount{
		ExternalAccount: pefi.ExternalAccount{
			Name:        name,
			OwnerId:     owner,
			Description: description,
		},
		Amount: pefi.MonetaryAmount{0, "SEK"},
	}, nil
}

func (s AccountService) Update(name string, new interface{}) error {
	return nil
}

func (s AccountService) Delete(name string) error {
	return nil
}

func (s AccountService) Get(name string) (interface{}, error) {
	return nil, errors.New("no such account")
}

func (s AccountService) Transfer(sender, receiver string) (string, error) {
	return "transfer id", nil
}

func (s AccountService) Deposit(name string, amount uint64) (string, error) {
	return "deposit id", nil
}

func (s AccountService) Withdraw(name string, amount uint64) (string, error) {
	return "withdraw id", nil
}

func SaveExternal(a pefi.ExternalAccount) error {
	//_, err := s.db.Exec("INSERT INTO external_accounts(name, owner-id, description) VALUES($1, $2, $3)", a.Name, a.OwnerId, a.Description)
	return nil
}

func SaveInternal(a pefi.InternalAccount) error {
	//_, err := s.db.Exec("INSERT INTO internal_accounts(name, owner-id, description, amount) VALUES($1, $2, $3, $4)", a.Name, a.OwnerId, a.Description, a.Amount)
	return nil
}
