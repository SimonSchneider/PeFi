package main

import (
	"context"
	"github.com/simonschneider/pefi/models"
)

type (
	label struct {
		extension string
		cl        *client
	}
)

func (l label) Extension() string {
	return l.extension
}

func (l label) GetNew() interface{} {
	return new(models.Label)
}

func (l label) GetAll(ctx context.Context) (allModels interface{}, err error) {
	//user := ctx.Value(userID).(int64)
	return nil, nil
}

func (l label) Get(ctx context.Context, id int64) (reqModel interface{}, err error) {
	return nil, nil
}

func (l label) Add(ctx context.Context, modModel interface{}) (err error) {
	return nil
}

func (l label) Del(ctx context.Context, id int64) (err error) {
	return nil
}

func (l label) Mod(ctx context.Context, id int64, modModel interface{}) (err error) {
	return nil
}
