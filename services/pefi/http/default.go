package http

import (
	"github.com/gorilla/mux"
	"github.com/simonschneider/pefi/middleware"
	"net/http"
)

type (
	Handler interface {
		Attach(r *mux.Router)
	}
)

func Init() *mux.Router {
	return mux.NewRouter()
}

func AttachAndStart(handlers ...Handler) {
	router := mux.NewRouter()

	for _, handler := range handlers {
		handler.Attach(router)
	}

	http.ListenAndServe(":8080", middleware.ApplyMiddleware(router,
		middleware.Json,
		middleware.General("accountService"),
		middleware.Timer,
		middleware.Context,
	))
}
