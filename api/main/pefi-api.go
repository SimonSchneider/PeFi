package main

import (
	"net/http"
	"pefi/api"
	//"pefi/logger"
	"pefi/router"
)

func main() {
	router := router.NewRouter(api.Routes)

	http.ListenAndServe(":22400", router)
}
