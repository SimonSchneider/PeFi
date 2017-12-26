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
	UserService struct {
		db *sql.DB
	}
)

func NewUserService(config *Config) (*UserService, error) {
	conn := getConnectionString(config)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &UserService{
		db,
	}, nil
}

func (s UserService) Create(ctx context.Context, name string) (*pefi.User, error) {
	user := &pefi.User{
		ID:   pefi.ID(uuid.NewV4().String()),
		Name: name,
	}
	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO users(id, name) VALUES($1, $2)")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	_, err = stmt.Exec(user.ID, user.Name)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return user, nil
}

func (s UserService) Update(ctx context.Context, id pefi.ID, new interface{}) error {
	return nil
}

func (s UserService) Delete(ctx context.Context, id pefi.ID) error {
	return nil
}

func (s UserService) Get(ctx context.Context, id pefi.ID) (*pefi.User, error) {
	var user pefi.User
	err := s.db.QueryRowContext(ctx, "SELECT id, name FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Name)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &user, err
}
