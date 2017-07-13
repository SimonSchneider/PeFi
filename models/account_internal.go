package models

type (
	InternalAccount struct {
		ExternalAccount
		Balance float64 `json:"balance" db:"balance"`
	}
)
