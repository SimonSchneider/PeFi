package models

import (
	"time"
)

type (
	Transaction struct {
		ID         int64     `json:"id" db:"id"`
		Time       time.Time `json:"time" db:"time"`
		Amount     float64   `json:"amount,number" db:"amount"`
		SenderID   int64     `json:"sender_id" db:"sender_id"`
		ReceiverID int64     `json:"receiver_id" db:"receiver_id"`
		LabelID    int64     `json:"label_id" db:"label_id"`
	}
)
