package router

import (
	"github.com/gorilla/mux"
	"net/http"
)

type (
	Adaptor func(http.HandlerFunc) http.HandlerFunc

	Route struct {
		Name        string
		Method      string
		Pattern     string
		HandlerFunc http.HandlerFunc
	}

	Routes []Route
)

//Adds the Adaptors to the http.Handler
func Handler(h http.HandlerFunc, adaptors ...Adaptor) http.HandlerFunc {
	for i := len(adaptors) - 1; i >= 0; i-- {
		h = adaptors[i](h)
	}
	return h
}

func NewRouter(routes Routes) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}
