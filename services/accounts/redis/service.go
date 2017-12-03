package redis

import (
	"fmt"
	"github.com/simonschneider/pefi/services/accounts"
)

type (
	Service struct {
		cache   map[string]interface{}
		service accounts.Service
	}
)

func NewService(next accounts.Service) *Service {
	return &Service{
		cache:   make(map[string]interface{}),
		service: next,
	}
}

func (s Service) OpenExternal(name, owner, description string) (*accounts.ExternalAccount, error) {
	return s.service.OpenExternal(name, owner, description)
}

func (s Service) OpenInternal(name, owner, description string) (*accounts.InternalAccount, error) {
	a, err := s.service.OpenInternal(name, owner, description)
	if err != nil {
		return a, err
	}
	s.cache[a.Name] = a
	fmt.Println("added to cache")
	return a, err
}

func (s Service) Update(name string, new interface{}) error {
	return s.service.Update(name, new)
}

func (s Service) Delete(name string) error {
	return s.service.Delete(name)
}

func (s Service) Get(name string) (interface{}, error) {
	if acc := s.cache[name]; acc != nil {
		return acc, nil
	}
	return s.service.Get(name)
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
