package pefi

import (
	"context"
)

type (
	UserService interface {
		Create(ctx context.Context, name string) (*User, error)
		Update(ctx context.Context, id ID, new interface{}) error
		Delete(ctx context.Context, id ID) error
		Get(ctx context.Context, id ID) (*User, error)
	}

	User struct {
		ID   ID     `json:id`
		Name string `json:name`
	}
)
