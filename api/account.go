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
	json.NewEncoder(w).Encode(a.Id)
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
}

func GetInternalAccount(w http.ResponseWriter, r *http.Request) {
}

func GetExternalAccount(w http.ResponseWriter, r *http.Request) {
}

func GetAllAccounts(w http.ResponseWriter, r *http.Request) {
	as, err := model.GetAllAccounts()
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(as)
}

func GetAccount(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	tmp := []string{}
	ids := []int64{}
	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		return
	}
	for _, id := range tmp {
		idi, _ := strconv.Atoi(id)
		ids = append(ids, int64(idi))
	}
	fmt.Println(ids)
	a, err := model.GetAccounts(ids)
	fmt.Println(a)
	fmt.Println(err)
	if err != nil {
		return
	}
	json.NewEncoder(w).Encode(a)
}

func DelAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accID, err := strconv.Atoi(vars["accountId"])
	if err != nil {
		return
	}
	err = model.DelAccount(int64(accID))
	if err != nil {
		return
	}
	return
}
