package database

import (
	"database/sql"
	"fmt"
	//This is ok because its a driver
	_ "github.com/lib/pq"
)

type (
	Client struct {
		info string
		db   *sql.DB
	}
)

//DB returns the DB object of the client
func (c *Client) DB() (db *sql.DB, err error) {
	if c.db.Ping() != nil {
		c.db, err = sql.Open("postgres", c.info)
		if err != nil {
			return nil, err
		}
		if err = c.db.Ping(); err != nil {
			return nil, err
		}
	}
	return c.db, nil
}

//NewClient return a new client
func NewClient(host string, port int, user string, password string, dbname string) (c *Client, err error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		//"password=%s dbname=%s sslmode=disable",
		"dbname=%s sslmode=disable",
		host, port, user, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	c = &Client{
		info: psqlInfo,
		db:   db,
	}
	return c, nil
}
