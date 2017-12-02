package main

import (
	"github.com/gorilla/mux"
	"github.com/simonschneider/pefi/middleware"
	"log"
	"net/http"
	"os"
)

func main() {
	logger := log.New(os.Stderr, "accounts ", log.Ldate|log.Ltime)
	logger.Println("starting")
	accountService := NewAccountService(&AccountServiceConfiguration{})
	accountController := NewAccountController(accountService)
	router := mux.NewRouter()
	accountRouter := router.PathPrefix("/api/accounts/").Subrouter()
	accountRouter.HandleFunc("/open/{type}", accountController.Open())

	log.Fatal(http.ListenAndServe(":8080", middleware.ApplyMiddleware(router,
		middleware.JsonMW,
		middleware.GeneralMW("accountService"),
		middleware.TimerMW,
		middleware.ContextMW,
	)))
}
