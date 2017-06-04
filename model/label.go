package model

import (
//"encoding/json"
)

type (
	Label struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
)
