package model

import (
	"encoding/json"
	"fmt"
	"pefi/model/redis"
	"strconv"
	"strings"
	"time"
)

type (
	Transactions []Transaction

	Transaction struct {
		Id         int64     `json:"id"`
		Time       time.Time `json:"time"`
		Amount     float64   `json:"amount,number"`
		SenderId   int64     `json:"sender_id"`
		ReceiverId int64     `json:"receiver_id"`
		LabelIds   []int64   `json:"label_ids"`
	}

	Loan struct {
		Transaction
		PaybackIds []int64 `json:"payback_ids"`
	}
)

var (
	transactionHeader = []string{
		"id",
		"time",
		"amount",
		"sender",
		"receiver",
		"labels",
	}
)

func (ts *Transactions) Header() (s []string) {
	return transactionHeader
}

func (ts *Transactions) Body() (s [][]string) {
	for _, t := range *ts {
		s = append(s, t.Table())
	}
	return s
}

func (ts *Transaction) Header() (s []string) {
	return transactionHeader
}

func (t *Transaction) Body() (s [][]string) {
	return [][]string{t.Table()}
}

func (t *Transaction) Table() (s []string) {
	s = []string{
		strconv.Itoa(int(t.Id)),
		t.Time.Format("2006-01-02"),
		fmt.Sprintf("%.2f", t.Amount),
		strconv.Itoa(int(t.SenderId)),
		strconv.Itoa(int(t.ReceiverId)),
	}
	labelIds := []string{}
	for _, id := range t.LabelIds {
		labelIds = append(labelIds, strconv.Itoa(int(id)))
	}
	s = append(s, strings.Join(labelIds, ","))
	return s
}

func (l *Loan) Table() (s []string) {
	s = l.Transaction.Table()
	paybackIds := []string{}
	for _, id := range l.PaybackIds {
		paybackIds = append(paybackIds, strconv.Itoa(int(id)))
	}
	s = append(s, strings.Join(paybackIds, ","))
	return s
}

func NewTransaction(t Transaction) (*Transaction, error) {
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

func GetTransactions() (ts []Transaction, err error) {
	vals, err := redis.HGetAll("Transaction")
	if err != nil {
		return
	}
	for _, val := range vals {
		t := new(Transaction)
		if err = json.Unmarshal([]byte(val), t); err != nil {
			return
		}
		ts = append(ts, *t)
	}
	return
}

func GetTransaction(id int64) (nt *Transaction, err error) {
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

//func (t *Transaction) Commit() error {
//err := t.Sender.send(t.Amount)
//if err != nil {
//return err
//}
//t.Receiver.receive(t.Amount)
//commit(t.Sender)
//commit(t.Receiver)
//return nil
//}

//func (t *Transaction) MarshalJSON() (data []byte, err error) {
//return json.Marshal(struct {
//Transaction
//Receiver_id int64 `json:"receiver_id"`
//Sender_id   int64 `json:"sender_id"`
//}{
//Transaction: *t,
//Receiver_id: t.Receiver.GetId(),
//Sender_id:   t.Sender.GetId(),
//})
//}
//func (t *Transaction) UnmarshalJSON(data []byte) (err error) {
//type _Transaction Transaction
//tmp := &struct {
//Sender_id   int64 `json:"sender_id"`
//Receiver_id int64 `json:"receiver_id"`
//*_Transaction
//}{
//_Transaction: (*_Transaction)(t),
//}
//if err = json.Unmarshal(data, tmp); err != nil {
//return
//}
//t.Sender, err = GetAccount(tmp.Sender_id)
//if err != nil {
//return
//}
//t.Receiver, err = GetAccount(tmp.Receiver_id)
//if err != nil {
//return
//}
//return
//}

//func CreateTransaction(data []byte) (t *Transaction, err error) {
//t = new(Transaction)
//if err = json.Unmarshal(data, t); err != nil {
//return
//}
//id, err := redis.HIncrBy("unique_ids", "Transaction", 1)
//if err != nil {
//return
//}
//t.Id = id
//output, err := json.Marshal(t)
//redis.HSet("Transaction", string(id), string(output))
//return
//}

//func GetTransaction(id int64) (t *Transaction, err error) {
//data, err := redis.HGet("Transaction", string(id))
//if err == nil {
//t = new(Transaction)
//err = json.Unmarshal([]byte(data), t)
//return
//}
//return
//}
