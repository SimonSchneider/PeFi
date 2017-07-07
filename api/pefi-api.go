package main

import (
	"net/http"
	"pefi/router"
)

func main() {
	router := router.NewRouter(Routes)

	http.ListenAndServe(":22400", router)
}
