package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"pefi/model/redis"
	"strconv"
	"time"
)

type (
	Transaction struct {
		Id         int64     `json:"id"`
		Time       time.Time `json:"time"`
		Amount     float64   `json:"amount,number"`
		SenderId   int64     `json:"sender_id"`
		ReceiverId int64     `json:"receiver_id"`
		LabelIds   []int64   `json:"label_ids"`
	}
)

func NewTransaction(in interface{}) (interface{}, error) {
	t, ok := in.(*Transaction)
	if !ok {
		return nil, errors.New("couldnt cast")
	}
	id, err := redis.HIncrBy("unique_ids", "Transaction", 1)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	t.Id = id
	jt, err := json.Marshal(t)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	redis.HSet("Transaction", strconv.Itoa(int(t.Id)), string(jt))
	//Commit Transaction
	return &t, err
}

func GetTransactions() (interface{}, error) {
	vals, err := redis.HGetAll("Transaction")
	if err != nil {
		return nil, err
	}
	var ts []Transaction
	for _, val := range vals {
		t := new(Transaction)
		if err = json.Unmarshal([]byte(val), t); err != nil {
			return nil, err
		}
		ts = append(ts, *t)
	}
	return &ts, err
}

func GetTransaction(id int64) (nt interface{}, err error) {
	val, err := redis.HGet("Transaction", strconv.Itoa(int(id)))
	if err != nil {
		fmt.Println(err)
		return
	}
	nt = new(Transaction)
	err = json.Unmarshal([]byte(val), nt)
	return
}

func DelTransaction(id int64) (err error) {
	err = redis.HDel("Transaction", strconv.Itoa(int(id)))
	if err != nil {
		fmt.Println(err)
	}
	return
}
