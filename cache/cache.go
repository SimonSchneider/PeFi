package cache

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	//"os"
)

type (
	//Client is the client object
	Client struct {
		address string
		db      int64
		conn    redis.Conn
	}
)

//Do is a wrapper of redigo conn.Do
func (c *Client) Do(command string, args ...interface{}) (reply interface{}, err error) {
	return c.conn.Do(command, args...)
}

const (
	lockedDB = 0
)

// NewCache returns a new cache object if it is able to connect
// to the it
func NewCache(host string, port string) (*Client, error) {
	address := host + ":" + port
	conn, err := redis.Dial("tcp", address, redis.DialDatabase(lockedDB))
	if err != nil {
		return nil, err
	}
	db, err := redis.Int64(conn.Do("INCR", "caches"))
	if err != nil {
		return nil, err
	}
	if db == lockedDB {
		return nil, errors.New("got unexpected response from redis. Is used for something else?")
	}
	return &Client{
		address: address,
		db:      db,
		conn:    conn,
	}, nil
}

//func getClient() (redis.Conn, error) {
//if conn != nil {
//return conn, nil
//}
//host := os.Getenv("redis-host")
//port := os.Getenv("redis-port")
//conn, err := redis.Dial("tcp", host+":"+port, redis.DialDatabase(1))
//return conn, err
//}
