package main

import (
	"database/sql"
)

type (
	MonetaryAmount struct {
		Amount   int64  `json:amount`
		Currency string `json:currency`
	}

	Account interface {
		Deposit(MonetaryAmount) error
		Withdraw(MonetaryAmount) error
	}

	ExternalAccount struct {
		Name        string `json:name`
		OwnerId     string `json:owner-id`
		Description string `json:description`
		storage     AccountStorage
	}

	InternalAccount struct {
		ExternalAccount
		Amount MonetaryAmount `json:amount`
	}

	AccountStorage interface {
		Save(Account)
		Load(Account)
	}

	SqlAccountStorage struct {
		db sql.DB
	}
)

func (a *ExternalAccount) Deposit(m MonetaryAmount) error {
	return nil
}

func (a *ExternalAccount) Withdraw(m MonetaryAmount) error {
	return nil
}

func (s *SqlAccountStorage) SaveExternal(a ExternalAccount) error {
	_, err := s.db.Exec("INSERT INTO external_accounts(name, owner-id, description) VALUES($1, $2, $3)", a.Name, a.OwnerId, a.Description)
	return err
}

func (s *SqlAccountStorage) SaveInternal(a InternalAccount) error {
	_, err := s.db.Exec("INSERT INTO internal_accounts(name, owner-id, description, amount) VALUES($1, $2, $3, $4)", a.Name, a.OwnerId, a.Description, a.Amount)
	return err
}
