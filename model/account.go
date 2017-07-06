package model

import (
	"encoding/json"
	"fmt"
	"io"
	"pefi/model/redis"
	"strconv"
	"strings"
)

type (
	Account interface {
		Table() []string
	}

	ExternalAccount struct {
		Id          int64   `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		LabelIds    []int64 `json:"label_ids"`
	}

	InternalAccount struct {
		ExternalAccount
		Balance float64 `json:"balance"`
	}
)

func (a *InternalAccount) Table() (s []string) {
	s = a.ExternalAccount.Table()
	s = append(s, fmt.Sprintf("%.2f", a.Balance))
	return s
}

func (a *ExternalAccount) Table() (s []string) {
	s = []string{
		strconv.Itoa(int(a.Id)),
		a.Name,
		a.Description,
	}
	labelIds := []string{}
	for _, id := range a.LabelIds {
		labelIds = append(labelIds, strconv.Itoa(int(id)))
	}
	s = append(s, strings.Join(labelIds, ","))
	return s
}

func (a *InternalAccount) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		InternalAccount
	}{
		InternalAccount: *a,
	})
}
func (a *ExternalAccount) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ExternalAccount
	}{
		ExternalAccount: *a,
	})
}

func GetExternalAccounts() (accs []ExternalAccount, err error) {
	vals, err := redis.HGetAll("ExternalAccount")
	if err != nil {
		return
	}
	for _, val := range vals {
		a := new(ExternalAccount)
		if err = json.Unmarshal([]byte(val), a); err != nil {
			return
		}
		accs = append(accs, *a)
	}
	return
}

func GetInternalAccounts() (accs []InternalAccount, err error) {
	vals, err := redis.HGetAll("InternalAccount")
	if err != nil {
		return
	}
	for _, val := range vals {
		a := new(InternalAccount)
		if err = json.Unmarshal([]byte(val), a); err != nil {
			return
		}
		accs = append(accs, *a)
	}
	return
}

func GetExternalAccount(id int64) (acc *ExternalAccount, err error) {
	val, err := redis.HGet("ExternalAccount", strconv.Itoa(int(id)))
	if err != nil {
		fmt.Println(err)
		return
	}
	acc = new(ExternalAccount)
	err = json.Unmarshal([]byte(val), acc)
	return
}

func GetInternalAccount(id int64) (acc *InternalAccount, err error) {
	val, err := redis.HGet("InternalAccount", strconv.Itoa(int(id)))
	if err != nil {
		fmt.Println(err)
		return
	}
	acc = new(InternalAccount)
	err = json.Unmarshal([]byte(val), acc)
	return
}

func DelExternalAccount(id int64) (err error) {
	err = redis.HDel("ExternalAccount", strconv.Itoa(int(id)))
	if err != nil {
		fmt.Println(err)
	}
	return
}

func DelInternalAccount(id int64) (err error) {
	err = redis.HDel("InternalAccount", strconv.Itoa(int(id)))
	if err != nil {
		fmt.Println(err)
	}
	return
}

//NewExternalAccount create a new External Account and return it
func NewExternalAccount(data io.Reader) (a *ExternalAccount, err error) {
	a = new(ExternalAccount)
	id, err := redis.HIncrBy("unique_ids", "Account", 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = json.NewDecoder(data).Decode(a); err != nil {
		fmt.Println(err)
		return
	}
	a.Id = id
	ma, err := json.Marshal(a)
	if err != nil {
		return
	}
	redis.HSet("ExternalAccount", strconv.Itoa(int(a.Id)), string(ma))
	return
}

//NewInternalAccount create a new External Account and return it
func NewInternalAccount(data io.Reader) (a *InternalAccount, err error) {
	a = new(InternalAccount)
	id, err := redis.HIncrBy("unique_ids", "Account", 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = json.NewDecoder(data).Decode(a); err != nil {
		fmt.Println(err)
		return
	}
	a.Id = id
	ma, err := json.Marshal(a)
	if err != nil {
		return
	}
	redis.HSet("InternalAccount", strconv.Itoa(int(a.Id)), string(ma))
	return
}
