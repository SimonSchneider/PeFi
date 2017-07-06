package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"pefi/model"
	"strconv"
)

//AddExternalAccount add a new external Account
func AddExternalAccount(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	a, err := model.NewExternalAccount(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(a)
}

//AddInternalAccount add a new external Account
func AddInternalAccount(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	a, err := model.NewInternalAccount(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(a)
}

func GetExternalAccounts(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	accs, err := model.GetExternalAccounts()
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(accs)
}

func GetInternalAccounts(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	accs, err := model.GetInternalAccounts()
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(accs)
}

func GetExternalAccount(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["accountId"])
	if err != nil {
		fmt.Println(err)
		return
	}
	acc, err := model.GetExternalAccount(int64(id))
	json.NewEncoder(w).Encode(acc)
}

func GetInternalAccount(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["accountId"])
	if err != nil {
		fmt.Println(err)
		return
	}
	acc, err := model.GetInternalAccount(int64(id))
	json.NewEncoder(w).Encode(acc)
}

func DelExternalAccount(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["accountId"])
	if err != nil {
		fmt.Println(err)
		return
	}
	err = model.DelExternalAccount(int64(id))
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func DelInternalAccount(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["accountId"])
	if err != nil {
		fmt.Println(err)
		return
	}
	err = model.DelInternalAccount(int64(id))
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}
