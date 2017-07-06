package api

import (
//"encoding/json"
//"fmt"
//"io/ioutil"
//"net/http"
//"pefi/model"
)

//func AddTransaction(w http.ResponseWriter, r *http.Request) {
//if r.Body == nil {
//http.Error(w, "no post", 400)
//fmt.Println("no post")
//return
//}
//b, err := ioutil.ReadAll(r.Body)
//if err != nil {
//fmt.Println(err)
//http.Error(w, err.Error(), 300)
//return
//}
//t, err := model.CreateTransaction(b)
//if err != nil {
//fmt.Println(err)
//http.Error(w, err.Error(), 300)
//return
//}
//if err := t.Commit(); err != nil {
//fmt.Println(err)
//http.Error(w, err.Error(), 300)
//return
//}
//json.NewEncoder(w).Encode(t)
//}
