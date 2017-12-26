package postgres

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"github.com/simonschneider/pefi"
)

type (
	AccountService struct {
		db *sql.DB
	}
)

func NewAccountService(config *Config) (*AccountService, error) {
	conn := getConnectionString(config)
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

func (s AccountService) Open(ctx context.Context, name string, ownerID pefi.ID, description string) (*pefi.Account, error) {
	acc := &pefi.Account{
		ID:          pefi.ID(uuid.NewV4().String()),
		Name:        name,
		OwnerID:     ownerID,
		Description: description,
		Balance:     pefi.MonetaryAmount{0, "SEK"},
	}
	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO accounts(id, name, description, owner_id, balance, currency) VALUES($1, $2, $3, $4, $5, $6)")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	_, err = stmt.Exec(acc.ID, acc.Name, acc.Description, acc.OwnerID, acc.Balance.Amount, acc.Balance.Currency)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return acc, nil
}

func (s AccountService) Update(ctx context.Context, id pefi.ID, new interface{}) error {
	return nil
}

func (s AccountService) Delete(ctx context.Context, id pefi.ID) error {
	return nil
}

func (s AccountService) Get(ctx context.Context, id pefi.ID) (*pefi.Account, error) {
	var acc pefi.Account
	err := s.db.QueryRowContext(ctx, "SELECT id, name, owner_id, description, balance, currency FROM accounts WHERE id = $1", id).
		Scan(&acc.ID, &acc.Name, &acc.OwnerID, &acc.Description, &acc.Balance.Amount, &acc.Balance.Currency)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &acc, err
}

func (s AccountService) GetAll(ctx context.Context, userID pefi.ID) ([]*pefi.Account, error) {
	var accs []*pefi.Account
	rows, err := s.db.QueryContext(ctx, "SELECT id, name, owner_id, description, balance, currency FROM accounts WHERE owner_id = $1", userID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var acc pefi.Account
		err := rows.Scan(&acc.ID, &acc.Name, &acc.OwnerID, &acc.Description, &acc.Balance.Amount, &acc.Balance.Currency)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		accs = append(accs, &acc)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return accs, err
}

func (s AccountService) Transfer(ctx context.Context, sender, receiver pefi.ID) (string, error) {
	return "transfer id", nil
}

func (s AccountService) Deposit(ctx context.Context, id pefi.ID, amount uint64) (string, error) {
	return "deposit id", nil
}

func (s AccountService) Withdraw(ctx context.Context, id pefi.ID, amount uint64) (string, error) {
	return "withdraw id", nil
}
