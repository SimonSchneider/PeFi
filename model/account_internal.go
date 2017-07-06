package model

import (
	"encoding/json"
	"fmt"
	"io"
	"pefi/model/redis"
	"strconv"
)

type (
	InternalAccounts []InternalAccount

	InternalAccount struct {
		ExternalAccount
		Balance float64 `json:"balance"`
	}
)

var (
	internalAccountHeader = append(
		externalAccountHeader,
		"balance",
	)
)

func (is *InternalAccounts) Header() (s []string) {
	return internalAccountHeader
}

func (is *InternalAccounts) Body() (s [][]string) {
	for _, i := range *is {
		s = append(s, i.Table())
	}
	return s
}

func (is *InternalAccounts) Footer() (s []string) {
	sum := 0.0
	for _, i := range *is {
		sum += i.Balance
	}
	for i := 0; i < len(internalAccountHeader); i++ {
		s = append(s, "")
	}
	s[len(s)-1] = fmt.Sprintf("%.2f", sum)
	s[len(s)-2] = "Total"
	return s
}

func (i *InternalAccount) Header() (s []string) {
	return internalAccountHeader
}

func (i *InternalAccount) Body() (s [][]string) {
	return [][]string{i.Table()}
}

func (i *InternalAccount) Footer() (s []string) {
	return []string{}
}

func (a *InternalAccount) Table() (s []string) {
	s = a.ExternalAccount.Table()
	s = append(s, fmt.Sprintf("%.2f", a.Balance))
	return s
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

func DelInternalAccount(id int64) (err error) {
	err = redis.HDel("InternalAccount", strconv.Itoa(int(id)))
	if err != nil {
		fmt.Println(err)
	}
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
