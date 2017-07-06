package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"pefi/model"
	"strconv"
)

func AddLabel(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	lab := new(model.Label)
	if err := json.NewDecoder(r.Body).Decode(lab); err != nil {
		fmt.Println(err)
		return
	}
	nlab, err := model.NewLabel(*lab)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(nlab)
}

func GetLabels(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	labs, err := model.GetLabels()
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(labs)
}

func GetLabel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["labelId"])
	if err != nil {
		fmt.Println(err)
		return
	}
	lab, err := model.GetLabel(int64(id))
	json.NewEncoder(w).Encode(lab)
}

func DelLabel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["labelId"])
	if err != nil {
		fmt.Println(err)
		return
	}
	err = model.DelLabel(int64(id))
	if err != nil {
		fmt.Println(err)
	}
	return
}
