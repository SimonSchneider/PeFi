package router

import (
	"github.com/gorilla/mux"
	"net/http"
)

type (
	//HttpDecorator func(...interface{}) func(http.HandlerFunc) http.HandlerFunc
	HTTPDecorator func(http.HandlerFunc) http.HandlerFunc

	Route struct {
		Name        string
		Method      string
		Pattern     string
		HandlerFunc http.HandlerFunc
	}

	Routes []Route
)

func Handler(h http.HandlerFunc, decorators ...HTTPDecorator) http.HandlerFunc {
	for i := len(decorators) - 1; i >= 0; i-- {
		h = decorators[i](h)
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
