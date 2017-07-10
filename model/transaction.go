package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"pefi/api/database"
	"pefi/model/redis"
	"strconv"
	"time"
)

type (
	Transaction struct {
		ID         int64     `json:"id"`
		Time       time.Time `json:"time"`
		Amount     float64   `json:"amount,number"`
		SenderID   int64     `json:"sender_id"`
		ReceiverID int64     `json:"receiver_id"`
		LabelID    int64     `json:"label_id"`
	}
)

func NewTransaction(c *database.Client) func(in interface{}) (interface{}, error) {
	return func(in interface{}) (interface{}, error) {
		t, ok := in.(*Transaction)

		//db, err := c.DB()
		//if err != nil {
		//return nil, err
		//}
		//rows, err := db.Query("SELECT * FROM users")
		//if err != nil {
		//fmt.Println("error")
		//}
		//_, err := db.Exec("INSERT INTO TRANSACTIONS")
		//defer rows.Close()
		//for rows.Next() {
		//var r string
		//if err = rows.Scan(&r); err != nil {
		//fmt.Println("error reading")
		//}
		//fmt.Printf("got row: %s\n", r)
		//}
		//if rows.Err(); err != nil {
		//fmt.Println("error reading 2")
		//}
		fmt.Println("fell through")

		if !ok {
			return nil, errors.New("couldnt cast")
		}
		id, err := redis.HIncrBy("unique_ids", "Transaction", 1)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		t.ID = id
		jt, err := json.Marshal(t)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		redis.HSet("Transaction", strconv.Itoa(int(t.ID)), string(jt))
		//Commit Transaction
		return &t, err
	}
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
