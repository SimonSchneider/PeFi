package main

import (
	"fmt"
	"net/http"
	"os"
	"pefi/api/database"
	"pefi/api/middleware"
	"pefi/cache"
)

func main() {
	cacheHost := os.Getenv("redis-host")
	cachePort := os.Getenv("redis-port")
	cache, err := cache.NewCache(cacheHost, cachePort)
	if err != nil {
		return
	}

	initMiddleware(cache)

	c, err := database.NewClient("postgres", 5432, "postgres", "", "pefi")
	if err != nil {
		fmt.Println(err)
	}

	router := middleware.NewRouter(getRoutes(c))

	http.ListenAndServe(":22400", router)
}
