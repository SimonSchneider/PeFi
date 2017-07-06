package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"pefi/model"
	"strconv"
)

func AddTransaction(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	t := new(model.Transaction)
	if err := json.NewDecoder(r.Body).Decode(t); err != nil {
		fmt.Println(err)
		return
	}
	nt, err := model.NewTransaction(*t)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(nt)
}

func GetTransactions(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ts, err := model.GetTransactions()
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(ts)
}

func GetTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["transactionId"])
	if err != nil {
		fmt.Println(err)
		return
	}
	t, err := model.GetTransaction(int64(id))
	json.NewEncoder(w).Encode(t)
}

func DelTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["transactionId"])
	if err != nil {
		fmt.Println(err)
		return
	}
	err = model.DelTransaction(int64(id))
	if err != nil {
		fmt.Println(err)
	}
	return
}
