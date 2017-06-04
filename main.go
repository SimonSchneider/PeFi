package main

import (
	"fmt"
	"html/template"
	"net/http"
	"pefi/api"
	"pefi/model"
)

var templates = template.Must(template.ParseFiles("static/daily.html"))

type (
	M map[string]interface{}

	Page struct {
		Title string
		Body  []byte
	}
)

func dailyHandel(w http.ResponseWriter, r *http.Request) {
	//id := r.URL.Path[len("/daily/"):]
	as, err := getAccounts()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = templates.ExecuteTemplate(w, "daily.html", M{
		"ExternalAccounts": as,
		"InternalAccounts": as,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func getAccounts() (as []model.Account, err error) {
	eacc_str := `
	{
		"name": "Ica",
		"description": "ica handlingskonto",
		"account_type": "external",
		"labels": [
		{
			"id":"lji3o",
			"name":"handling",
			"description":"stuff i buy"
		},
		{
			"id":"ljio23o",
			"name":"others",
			"description":"other stuff"
		}]
	}`
	iacc_str := `
	{
		"name": "Ica",
		"description": "ica handlingskonto",
		"account_type": "internal",
		"balance": 953.29,
		"labels": [
		{
			"id":"lji3o",
			"name":"handling",
			"description":"stuff i buy"
		},
		{
			"id":"ljio23o",
			"name":"others",
			"description":"other stuff"
		}]
	}`
	tmpe, err := model.CreateAccount([]byte(eacc_str))
	if err != nil {
		return
	}
	as = append(as, tmpe)
	tmpi, err := model.CreateAccount([]byte(iacc_str))
	if err != nil {
		return
	}
	as = append(as, tmpi)
	return
}

func getTransactions() (ts []*model.Transaction) {
	tran_str := `
	{
	"amount": 242.2,
	"time": "2012-04-23T18:25:43.511Z",
	"sender_id": 2,
	"receiver_id": 1
	}`
	t, err := model.CreateTransaction([]byte(tran_str))
	if err != nil {
		fmt.Println(err)
		return
	}
	ts = append(ts, t)
	return
}

func main() {
	http.HandleFunc("/daily", dailyHandel)
	//http.HandleFunc("/api/getia", getAccounts)
	http.HandleFunc("/api/add/transaction", api.AddTransaction)
	http.HandleFunc("/api/add/account", api.AddAccount)
	http.HandleFunc("/api/get/all/accounts", api.GetAllAccounts)

	as, err := getAccounts()
	if err != nil {
		fmt.Println(err)
	}

	for _, a := range as {
		tmp, err := model.GetAccount(a.GetId())
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(tmp)
	}

	ts := getTransactions()

	for _, t := range ts {
		tmp, err := model.GetTransaction(t.Id)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(*tmp)
		fmt.Println(tmp.Sender)
		fmt.Println(tmp.Receiver)
	}
	http.ListenAndServe(":8080", nil)
}
