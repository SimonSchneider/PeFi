package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"pefi/model/redis"
	"strconv"
)

type (
	Label struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		CategorieID int64  `json:"category_id"`
	}
)

func GetLabels() (interface{}, error) {
	vals, err := redis.HGetAll("Label")
	if err != nil {
		return nil, err
	}
	var labs []Label
	for _, val := range vals {
		l := new(Label)
		if err = json.Unmarshal([]byte(val), l); err != nil {
			return nil, err
		}
		labs = append(labs, *l)
	}
	return &labs, nil
}

func GetLabel(id int64) (lab interface{}, err error) {
	val, err := redis.HGet("Label", strconv.Itoa(int(id)))
	if err != nil {
		fmt.Println(err)
		return
	}
	lab = new(Label)
	err = json.Unmarshal([]byte(val), lab)
	return
}

func DelLabel(id int64) (err error) {
	err = redis.HDel("Label", strconv.Itoa(int(id)))
	if err != nil {
		fmt.Println(err)
	}
	return
}

func NewLabel(in interface{}) (nlab interface{}, err error) {
	lab, ok := in.(*Label)
	if !ok {
		return nil, errors.New("couldnt cast")
	}
	id, err := redis.HIncrBy("unique_ids", "Label", 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	lab.ID = id
	jlab, err := json.Marshal(lab)
	if err != nil {
		fmt.Println(err)
		return
	}
	redis.HSet("Label", strconv.Itoa(int(lab.ID)), string(jlab))
	return &lab, err
}
