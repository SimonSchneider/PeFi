package redis

import (
	"github.com/go-redis/redis"
	"os"
)

var client *redis.Client

func GetClient() (*redis.Client, error) {
	if client != nil {
		return client, nil
	}
	host := os.Getenv("redis-host")
	port := os.Getenv("redis-port")
	client = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: "",
		DB:       0,
	})
	_, err := client.Ping().Result()
	return client, err
}
