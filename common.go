package pefi

import (
	"errors"
	"strconv"
)

type (
	MonetaryAmount struct {
		Amount   int64  `json:amount`
		Currency string `json:currency`
	}

	ID string
)

func (_ ID) ImplementsGraphQLType(name string) bool {
	return name == "ID"
}

func (id *ID) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case string:
		*id = ID(input)
	default:
		err = errors.New("wrong type")
	}
	return err
}

func (id ID) MarshalJSON() ([]byte, error) {
	return strconv.AppendQuote(nil, string(id)), nil
}
