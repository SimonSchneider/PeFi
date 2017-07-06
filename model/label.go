package model

import (
	"encoding/json"
	"fmt"
	"io"
	"pefi/model/redis"
	"strconv"
)

type (
	Labels []Label

	Label struct {
		Id          int64  `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
)

var (
	labelHeader = []string{
		"id",
		"name",
		"desc",
	}
)

func (ls *Labels) Header() (s []string) {
	return labelHeader
}

func (ls *Labels) Body() (s [][]string) {
	for _, l := range *ls {
		s = append(s, l.Table())
	}
	return s
}

func (ls *Labels) Footer() (s []string) {
	return []string{}
}

func (l *Label) Header() (s []string) {
	return labelHeader
}

func (l *Label) Body() (s [][]string) {
	return [][]string{l.Table()}
}

func (l *Label) Footer() (s []string) {
	return []string{}
}

func (l *Label) Table() (s []string) {
	return []string{
		strconv.Itoa(int(l.Id)),
		l.Name,
		l.Description,
	}
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

func DelLabel(id int64) (err error) {
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

func DecodeLabel(r io.Reader) (label interface{}, err error) {
	label = new(Label)
	err = json.NewDecoder(r).Decode(label)
	return
}

func EncodeLabels(ls []Label, w io.Writer) (err error) {
	return json.NewEncoder(w).Encode(ls)
}

func EncodeLabel(l Label, w io.Writer) (err error) {
	return json.NewEncoder(w).Encode(l)
}
