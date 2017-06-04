package model

import (
	"encoding/json"
	"fmt"
	"pefi/model/redis"
	"strconv"
)

type (
	Account interface {
		GetId() int64
		setId(id int64)
		receive(amount float64)
		send(amount float64) error
	}

	ExternalAccount struct {
		Id          int64   `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Labels      []Label `json:"labels"`
	}

	InternalAccount struct {
		ExternalAccount
		Balance float64 `json:"balance"`
	}
)

func (a *ExternalAccount) setId(id int64) {
	a.Id = id
}
func (a *InternalAccount) setId(id int64) {
	a.Id = id
}

func (a *ExternalAccount) GetId() int64 {
	return a.Id
}
func (a *InternalAccount) GetId() int64 {
	return a.Id
}

func (a *ExternalAccount) receive(amount float64) {}
func (a *InternalAccount) receive(amount float64) {
	a.Balance += amount
}

func (a *ExternalAccount) send(amount float64) error {
	return nil
}
func (a *InternalAccount) send(amount float64) error {
	if amount > a.Balance {
		return fmt.Errorf("no coverage for %.02f on account %d with balance %.02f", amount, a.Id, a.Balance)
	}
	a.Balance -= amount
	return nil
}

func (a *InternalAccount) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		InternalAccount
		Account_type string `json:"account_type"`
	}{
		InternalAccount: *a,
		Account_type:    "internal",
	})
}
func (a *ExternalAccount) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ExternalAccount
		Account_type string `json:"account_type"`
	}{
		ExternalAccount: *a,
		Account_type:    "external",
	})
}

// Determines the account type of the json data string
func jsonToAccount(data []byte) (a Account, err error) {
	type extr struct {
		Account_type string `json:"account_type"`
	}
	t := new(extr)
	json.Unmarshal(data, t)
	switch t.Account_type {
	case "external":
		a = new(ExternalAccount)
		if err = json.Unmarshal([]byte(data), a); err == nil {
			return
		}
	case "internal":
		a = new(InternalAccount)
		if err = json.Unmarshal([]byte(data), a); err == nil {
			return
		}
	}
	err = fmt.Errorf("not valid account type")
	return
}

// Returns the account that has the ID, the account type will
// be decided dynamically
func GetAccount(id int64) (a Account, err error) {
	rc, err := redis.GetClient()
	data, err := rc.HGet("Account", string(id)).Result()
	if err != nil {
		return
	}
	a, err = jsonToAccount([]byte(data))
	return
}

// Creates an account from a json data string
// The account type will automatically be determined
// Will get a new Id from the database and give it to the account
func CreateAccount(data []byte) (a Account, err error) {
	a, err = jsonToAccount(data)
	if err != nil {
		return
	}
	rc, err := redis.GetClient()
	if err != nil {
		return
	}
	id, err := rc.HIncrBy("unique_ids", "Account", 1).Result()
	if err != nil {
		return
	}
	a.setId(id)
	output, err := json.Marshal(a)
	if err = rc.HSet("Account", string(id), output).Err(); err != nil {
		return
	}
	return
}

func GetAllAccounts() (as []Account, err error) {
	rc, err := redis.GetClient()
	if err != nil {
		return
	}
	ids, err := rc.HKeys("Account").Result()
	if err != nil {
		return
	}
	for _, sid := range ids {
		id, _ := strconv.Atoi(sid)
		a, _ := GetAccount(int64(id))
		as = append(as, a)
	}
	return
}
