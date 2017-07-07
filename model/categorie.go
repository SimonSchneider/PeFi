package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"pefi/model/redis"
	"strconv"
)

type (
	Categorie struct {
		Id          int64   `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		LabelIds    []int64 `json:"label_ids"`
		ChildrenIds []int64 `json:"children_ids"`
	}
)

func GetCategories() (interface{}, error) {
	vals, err := redis.HGetAll("Categorie")
	if err != nil {
		return nil, err
	}
	var cs []Categorie
	for _, val := range vals {
		c := new(Categorie)
		if err = json.Unmarshal([]byte(val), c); err != nil {
			return nil, err
		}
		cs = append(cs, *c)
	}
	return &cs, nil
}

func GetCategorie(id int64) (c interface{}, err error) {
	val, err := redis.HGet("Categorie", strconv.Itoa(int(id)))
	if err != nil {
		fmt.Println(err)
		return
	}
	c = new(Categorie)
	err = json.Unmarshal([]byte(val), c)
	return
}

func DelCategorie(id int64) (err error) {
	err = redis.HDel("Categorie", strconv.Itoa(int(id)))
	if err != nil {
		fmt.Println(err)
	}
	return
}

func NewCategorie(in interface{}) (nc interface{}, err error) {
	c, ok := in.(*Categorie)
	if !ok {
		return nil, errors.New("couldnt cast")
	}
	fmt.Println(c)
	fmt.Println("testing")
	id, err := redis.HIncrBy("unique_ids", "Categorie", 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Id = id
	jc, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	redis.HSet("Categorie", strconv.Itoa(int(c.Id)), string(jc))
	return &c, err
}
