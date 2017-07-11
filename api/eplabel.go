package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocraft/dbr"
	"github.com/gorilla/mux"
	"net/http"
	"pefi/api/models"
	"strconv"
)

type (
	labelT struct {
		models.Label
		User int64 `db:"user_id"`
	}
	labelEP struct {
	}
)

var (
	sess *dbr.Session
)

func Init(s string) {
	dbt, err := dbr.Open("postgres", s, nil)
	if err != nil {
		fmt.Println(err)
	}
	sess = dbt.NewSession(nil)
}

func (l *labelEP) URL() string {
	return "/labels"
}

func (l *labelEP) GetNewInstance() interface{} {
	return new(models.Label)
}

func Tmp() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := strconv.Atoi(r.Header.Get("user"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		e := labelEP{}
		r.ParseForm()
		models, err := e.GetAll(int64(user), r.Form)
		fmt.Println(err)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err = json.NewEncoder(w).Encode(models); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func (l *labelEP) GetAll(user int64, form map[string][]string) (interface{}, error) {
	q := sess.Select("label.id", "label.name", "label.description", "label.category_id").
		From("label").
		Join("category", "label.category_id=category.id").
		Where("category.user_id=?", user)
	var labels []models.Label
	orderBy := form["orderBy"]
	for _, order := range orderBy {
		if order != "" {
			q = q.OrderBy("label." + order)
		}
	}
	var limit uint64
	tmp := form["limit"]
	if len(tmp) != 0 {
		t, _ := strconv.Atoi(form["limit"][0])
		limit = uint64(t)
	}
	if limit != 0 {
		q = q.Limit(limit)
	}
	q.LoadStructs(&labels)
	fmt.Println(labels)
	return labels, nil
}

func (l *labelEP) Add(user int64, IDENT interface{}) error {
	c, ok := IDENT.(models.Label)
	if !ok {
		return errors.New("error in incoming data")
	}
	q := `
	INSERT INTO label
	(name, description, category_id)
	SELECT
	:name, :description, :category_id
	WHERE :user_id IN (
		SELECT user_id FROM category
		WHERE id=:category_id
	)
	`
	t := &labelT{c, user}
	_, err := db.NamedExec(q, t)
	return err
}

func TmpG() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		user, err := strconv.Atoi(r.Header.Get("user"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		lep := labelEP{}
		mod, err := lep.Get(int64(user), int64(id))
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				w.WriteHeader(http.StatusNotFound)
			default:
				w.WriteHeader(http.StatusBadRequest)
			}
			return
		}
		if err = json.NewEncoder(w).Encode(mod); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
	})
}

func (l *labelEP) Get(user, id int64) (interface{}, error) {
	q := `
	SELECT l.id, l.name, l.description, l.category_id
	FROM label AS l LEFT JOIN category AS c ON c.id = l.category_id
	WHERE c.user_id=$1 AND l.id=$2
	`
	out := models.Label{}
	err := db.Get(&out, q, user, id)
	return &out, err
}

func (l *labelEP) Del(user, id int64) error {
	return nil
}

func (l *labelEP) Mod(user, id int64, IDENT interface{}) error {

	return nil
}
