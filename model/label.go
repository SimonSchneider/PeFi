package model

import (
	"encoding/json"
	"fmt"
	"pefi/model/redis"
	"strconv"
)

type (
	Label struct {
		Id          int64  `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
)

func (l *Label) Table() (s []string) {
	s = []string{
		strconv.Itoa(int(l.Id)),
		l.Name,
		l.Description,
	}
	return s
}

func (l *Label) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Label
	}{
		Label: *l,
	})
}

func GetLabels() (labs []Label, err error) {
	vals, err := redis.HGetAll("Label")
	if err != nil {
		return
	}
	for _, val := range vals {
		l := new(Label)
		if err = json.Unmarshal([]byte(val), l); err != nil {
			return
		}
		labs = append(labs, *l)
	}
	return
}

func GetLabel(id int64) (lab *Label, err error) {
	val, err := redis.HGet("Label", strconv.Itoa(int(id)))
	if err != nil {
		fmt.Println(err)
		return
	}
	lab = new(Label)
	err = json.Unmarshal([]byte(val), lab)
	return
}

func DelLabel(id int64) (lab *Label, err error) {
	err = redis.HDel("Label", strconv.Itoa(int(id)))
	if err != nil {
		fmt.Println(err)
	}
	return
}

func NewLabel(lab Label) (nlab *Label, err error) {
	id, err := redis.HIncrBy("unique_ids", "Label", 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	lab.Id = id
	jlab, err := json.Marshal(lab)
	if err != nil {
		fmt.Println(err)
		return
	}
	redis.HSet("Label", strconv.Itoa(int(lab.Id)), string(jlab))
	return &lab, err
}
