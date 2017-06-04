package model

import (
	//"pefi/model/db"
	"encoding/json"
	"pefi/model/redis"
)

var root_node = root{}

type (
	root struct {
		Children []Categorie
	}

	Categorie struct {
		Id       int64       `json:"id"`
		Name     string      `json:"name"`
		Children []Categorie `json:"children"`
	}
)

func GetCategorie(id string) (c *Categorie, err error) {
	c = new(Categorie)
	rc, err := redis.GetClient()
	data, err := rc.Get("categorie:" + string(id)).Result()
	if err = json.Unmarshal([]byte(data), c); err != nil {
		return
	}
	return
}

func CreateCategorie(data string) (c *Categorie, err error) {
	c = new(Categorie)
	if err = json.Unmarshal([]byte(data), c); err != nil {
		return
	}
	rc, err := redis.GetClient()
	//add to database and redis
	//replace by database insert to get id
	id, err := rc.HIncrBy("unique_ids", "categorie", 1).Result()
	if err != nil {
		return
	}
	c.Id = id
	output, err := json.Marshal(c)
	if err != nil {
		return
	}
	if err = rc.Set("categorie:"+string(id), output, 0).Err(); err != nil {
		return
	}
	return
}
