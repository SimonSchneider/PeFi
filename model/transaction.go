package model

//import (
//"encoding/json"
//"pefi/model/redis"
//"time"
//)

//type (
//Transaction struct {
//Id       int64     `json:"id"`
//Amount   float64   `json:"amount,number"`
//Time     time.Time `json:"time"`
//Sender   Account   `json:"-"`
//Receiver Account   `json:"-"`
//Labels   []Label   `json:"labels"`
//}

//Loan struct {
//Transaction
//payback []Transaction
//}
//)

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
