package models

import (
	"github.com/jmoiron/sqlx"
)

type (
	// Endpoint is an interface for all models to implement
	// This will be called by the middleware for the given
	// endpoint
	Endpoint interface {
		URL() string
		Add(user int64) error
		GetAll(user int64) (interface{}, error)
		Get(user int64, id int64) (Endpoint, error)
		Del(user int64, id int64) error
		Mod(user int64, id int64) error
	}
)

var (
	db *sqlx.DB
)

func InitDB(newDB *sqlx.DB) {
	db = newDB
}
