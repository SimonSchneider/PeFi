package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"pefi/cache"
	"pefi/logger"
	"pefi/router"
	"strconv"
)

type (
	apiHandlers struct {
		gets http.HandlerFunc
		get  http.HandlerFunc
		add  http.HandlerFunc
		del  http.HandlerFunc
	}
	route struct {
		name     string
		handlers apiHandlers
	}

	addFunc  func(mod interface{}) (newMod interface{}, err error)
	getFunc  func(id int64) (newMod interface{}, err error)
	delFunc  func(id int64) (err error)
	getsFunc func() (newMods interface{}, err error)
)

func mwAdd(mod interface{}, apiFunc addFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(mod); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		newMod, err := apiFunc(mod)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err = json.NewEncoder(w).Encode(newMod); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func mwGet(apiFunc getFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		newMod, err := apiFunc(int64(id))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err = json.NewEncoder(w).Encode(newMod); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func mwGets(apiFunc getsFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newMods, err := apiFunc()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err = json.NewEncoder(w).Encode(newMods); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func mwDel(apiFunc delFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		err = apiFunc(int64(id))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func createRoutes(rs []route) router.Routes {
	var routerRoutes router.Routes
	for _, r := range rs {
		routerRoutes = append(routerRoutes,
			router.Route{
				"Get all " + r.name,
				"GET",
				"/" + r.name,
				router.Handler(r.handlers.gets,
					logger.HTTPLogger(r.name), cache.HTTPCache(120)),
			},
			router.Route{
				"Get id " + r.name,
				"GET",
				"/" + r.name + "/{id}",
				router.Handler(r.handlers.get,
					logger.HTTPLogger(r.name), cache.HTTPCache(120)),
			},
			router.Route{
				"Del id " + r.name,
				"DEL",
				"/" + r.name + "/{id}",
				router.Handler(r.handlers.del,
					logger.HTTPLogger(r.name), cache.HTTPWipe("/"+r.name)),
			},
			router.Route{
				"Add " + r.name,
				"POST",
				"/" + r.name,
				router.Handler(r.handlers.add,
					logger.HTTPLogger(r.name), cache.HTTPWipe("/"+r.name)),
			},
		)
	}
	return routerRoutes
}
