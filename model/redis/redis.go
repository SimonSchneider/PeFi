package redis

import (
	"github.com/garyburd/redigo/redis"
	"os"
)

var conn redis.Conn

func getClient() (redis.Conn, error) {
	if conn != nil {
		return conn, nil
	}
	host := os.Getenv("redis-host")
	port := os.Getenv("redis-port")
	conn, err := redis.Dial("tcp", host+":"+port)
	return conn, err
}

func HGetAll(key string) (val map[string]string, err error) {
	rc, err := getClient()
	if err != nil {
		return
	}
	return redis.StringMap(rc.Do("HGETALL", key))
}

func HSet(key string, hash string, val string) (err error) {
	rc, err := getClient()
	if err != nil {
		return
	}
	_, err = rc.Do("HSET", key, hash, val)
	return
}

func HDel(key string, hash string) (err error) {
	rc, err := getClient()
	if err != nil {
		return
	}
	_, err = rc.Do("HDEL", key, hash)
	return
}

func HIncrBy(key string, hash string, inc int64) (val int64, err error) {
	rc, err := getClient()
	if err != nil {
		return
	}
	return redis.Int64(rc.Do("HINCRBY", key, hash, inc))
}

func HGet(key string, hash string) (val string, err error) {
	rc, err := getClient()
	if err != nil {
		return
	}
	return redis.String(rc.Do("HGET", key, hash))
}
