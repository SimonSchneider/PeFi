package main

import (
	"net/http"
	"pefi/api"
)

func main() {
	//http.HandleFunc("/add/transaction", api.AddTransaction)
	//http.HandleFunc("/account/add", api.AddAccount)
	//http.HandleFunc("/account/get", api.GetAccount)
	//http.HandleFunc("/account/get/all", api.GetAllAccounts)

	router := api.NewRouter()

	http.ListenAndServe(":22400", router)
}
