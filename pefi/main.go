package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type (
	Actor interface {
		receive(amount float64)
		send(amount float64)
	}

	Account struct {
		Id          string  `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Labels      []Label `json:"labels"`
	}

	PersonalAccount struct {
		Balance float64 `json:"balance"`
	}

	Label struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	Transaction struct {
		Id       string    `json:"id"`
		Amount   float64   `json:"amount"`
		Date     time.Time `json:"time"`
		Sender   Actor     `json:"sender"`
		Receiver Actor     `json:"receiver"`
		Labels   []Label   `json:"labels"`
	}

	Loan struct {
		Transaction
	}
)

func (a *Account) receive(amount float64) {}
func (a *Account) send(amount float64)    {}

func (a *PersonalAccount) receive(amount float64) {
	a.Balance += amount
}
func (a *PersonalAccount) send(amount float64) {
	a.Balance -= amount
}

func main() {
	json_str := `
	{
		"id": "e190j2l",
		"name": "Ica",
		"description": "ica handlingskonto",
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
	a := new(Account)
	err := json.Unmarshal([]byte(json_str), a)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(a)

	t := Transaction{
		Id:       "helio3",
		Amount:   242.2,
		Date:     time.Now(),
		Sender:   a,
		Receiver: a,
		Labels:   a.Labels}
	fmt.Println(t)
}
