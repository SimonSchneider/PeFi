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
	ExternalAccounts []ExternalAccount

	ExternalAccount struct {
		Id          int64   `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		LabelIds    []int64 `json:"label_ids"`
	}
)

var (
	externalAccountHeader = []string{
		"id",
		"name",
		"description",
		"labels",
	}
)

func (es *ExternalAccounts) Header() (s []string) {
	return externalAccountHeader
}

func (es *ExternalAccounts) Body() (s [][]string) {
	for _, e := range *es {
		s = append(s, e.Table())
	}
	return s
}

func (es *ExternalAccounts) Footer() (s []string) {
	return []string{}
}

func (e *ExternalAccount) Header() (s []string) {
	return externalAccountHeader
}

func (e *ExternalAccount) Body() (s [][]string) {
	return [][]string{e.Table()}
}

func (e *ExternalAccount) Footer() (s []string) {
	return []string{}
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

func DelExternalAccount(id int64) (err error) {
	err = redis.HDel("ExternalAccount", strconv.Itoa(int(id)))
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
