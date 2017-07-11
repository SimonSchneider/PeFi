package models

type (
	Category struct {
		ID          int64  `json:"id" db:"id"`
		Name        string `json:"name" db:"name"`
		Description string `json:"description" db:"description"`
	}
)

// Name returns the name of category endpoint
func (c *Category) URL() string {
	return "/categories"
}

func (c *Category) GetAll(user int64) (interface{}, error) {
	query := `
		SELECT id, name, description
		FROM category
		WHERE user_id::int=$1::bigint`
	categories := []Category{}
	err := db.Select(&categories, query, user)
	return &categories, err
}

func (c *Category) Add(user int64) error {
	query := `
		INSERT INTO category (id, name, description, user_id)
		VALUES(default,:name,:description,:user_id)`
	type tmp struct {
		Category
		User int64 `db:"user_id"`
	}
	t := &tmp{*c, int64(user)}
	_, err := db.NamedExec(query, t)
	return err
}

func (c *Category) Get(user int64, id int64) (Endpoint, error) {
	query := `
		SELECT id, name, description
		FROM category
		WHERE user_id=$1::bigint AND id=$2::bigint`

	category := Category{}
	err := db.Get(&category, query, user, id)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (c *Category) Del(user int64, id int64) error {
	query := `
		DELETE FROM category WHERE
		id=:id AND user_id=:user_id
		`
	category := Category{ID: id}
	type tmp struct {
		Category
		User int64 `db:"user_id"`
	}
	t := &tmp{category, user}
	_, err := db.NamedExec(query, t)
	return err
}

func (c *Category) Mod(user int64, id int64) error {
	query := `
		UPDATE category SET
		name = :name,
		description = :description
		WHERE id=:id AND user_id=:user_id`
	type tmp struct {
		Category
		User int64 `db:"user_id"`
	}
	c.ID = id
	t := tmp{*c, user}
	_, err := db.NamedExec(query, t)
	return err
}
