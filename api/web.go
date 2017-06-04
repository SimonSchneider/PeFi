package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"pefi/model"
)

func AddAccount(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
	a, err := model.CreateAccount(b)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("account", a)
	return
}

func GetAllAccounts(w http.ResponseWriter, r *http.Request) {
	as, err := model.GetAllAccounts()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(as)
	json.NewEncoder(w).Encode(as)
}

func AddTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "no post", 400)
		fmt.Println("no post")
		return
	}
	t := new(model.Transaction)
	if err := json.NewDecoder(r.Body).Decode(t); err != nil {
		fmt.Println(err)
		return
	}
	if err := t.Commit(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(t)
}
