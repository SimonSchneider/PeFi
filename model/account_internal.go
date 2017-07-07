package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"pefi/model/redis"
	"strconv"
)

type (
	InternalAccount struct {
		ExternalAccount
		Balance float64 `json:"balance"`
	}
)

func GetInternalAccounts() (interface{}, error) {
	vals, err := redis.HGetAll("InternalAccount")
	if err != nil {
		return nil, err
	}
	var accs []InternalAccount
	for _, val := range vals {
		a := new(InternalAccount)
		if err = json.Unmarshal([]byte(val), a); err != nil {
			return nil, err
		}
		accs = append(accs, *a)
	}
	return &accs, nil
}

func GetInternalAccount(id int64) (acc interface{}, err error) {
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
func NewInternalAccount(in interface{}) (na interface{}, err error) {
	a, ok := in.(*InternalAccount)
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
	redis.HSet("InternalAccount", strconv.Itoa(int(a.Id)), string(ma))
	return &a, err
}
