package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/simonschneider/pefi/services/pefi"
)

type (
	AccountService struct {
		db *sql.DB
	}
)

func NewAccountService(config *Config) (*AccountService, error) {
	conn := getConnectionString(config)
	fmt.Println(conn)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &AccountService{
		db,
	}, nil
}

func (s AccountService) Open(ctx context.Context, name, owner, description string) (*pefi.Account, error) {
	return &pefi.Account{
		Name:        name,
		OwnerId:     owner,
		Description: description,
		Amount:      pefi.MonetaryAmount{0, "SEK"},
	}, nil
}

func (s AccountService) Update(ctx context.Context, name string, new interface{}) error {
	return nil
}

func (s AccountService) Delete(ctx context.Context, name string) error {
	return nil
}

func (s AccountService) Get(ctx context.Context, name string) (*pefi.Account, error) {
	return nil, errors.New("no such account")
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

func Save(a pefi.Account) error {
	//_, err := s.db.Exec("INSERT INTO internal_accounts(name, owner-id, description, amount) VALUES($1, $2, $3, $4)", a.Name, a.OwnerId, a.Description, a.Amount)
	return nil
}
