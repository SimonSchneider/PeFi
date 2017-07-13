package models

type (
	Label struct {
		ID          int64  `json:"id" db:"id"`
		Name        string `json:"name" db:"name"`
		Description string `json:"description" db:"description"`
		CategoryID  int64  `json:"category_id" db:"category_id"`
	}
)
