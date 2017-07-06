package model

type (
	Tabular interface {
		Header() []string
		Body() [][]string
	}
)
