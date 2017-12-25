package http

import (
	"github.com/gorilla/mux"
	"github.com/simonschneider/pefi/middleware"
	"log"
	"net/http"
	"net/url"
)

type (
	Handler interface {
		Attach(r *mux.Router)
	}
)

func GetRouter(handlers ...Handler) *mux.Router {
	router := mux.NewRouter()
	for _, handler := range handlers {
		handler.Attach(router)
	}
	return router
}

func AttachAndStart(handlers ...Handler) {
	router := GetRouter(handlers...)

	log.Fatal(http.ListenAndServe(":8080", middleware.ApplyMiddleware(router,
		middleware.Json,
		middleware.General("accountService"),
		middleware.Timer,
		//middleware.Context,
	)))
}

func getRequest(url *url.URL) (*http.Request, error) {
	url.Host = "localhost:8080"
	url.Scheme = "http"
	return http.NewRequest(
		"GET",
		url.String(),
		nil,
	)
}
