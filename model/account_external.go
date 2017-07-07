package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"pefi/model/redis"
	"strconv"
)

type (
	ExternalAccount struct {
		Id          int64   `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		LabelIds    []int64 `json:"label_ids"`
	}
)

func GetExternalAccounts() (interface{}, error) {
	vals, err := redis.HGetAll("ExternalAccount")
	if err != nil {
		return nil, err
	}
	var accs []ExternalAccount
	for _, val := range vals {
		a := new(ExternalAccount)
		if err = json.Unmarshal([]byte(val), a); err != nil {
			return nil, err
		}
		accs = append(accs, *a)
	}
	return &accs, nil
}

func GetExternalAccount(id int64) (acc interface{}, err error) {
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
func NewExternalAccount(in interface{}) (na interface{}, err error) {
	a, ok := in.(*ExternalAccount)
	if !ok {
		return nil, errors.New("couldnt cast")
	}
	id, err := redis.HIncrBy("unique_ids", "Account", 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	a.Id = id
	ma, err := json.Marshal(a)
	if err != nil {
		return
	}
	redis.HSet("ExternalAccount", strconv.Itoa(int(a.Id)), string(ma))
	return &a, err
}
