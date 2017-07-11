package models

type (
	Label struct {
		ID          int64  `json:"id" db:"id"`
		Name        string `json:"name" db:"name"`
		Description string `json:"description" db:"description"`
		CategoryID  int64  `json:"category_id" db:"category_id"`
	}
)

// Name returns the name of category endpoint
func (c *Label) URL() string {
	return "/labels"
}

func (c *Label) GetAll(user int64) (interface{}, error) {
	query := `
		SELECT id, name, description, category_id
		FROM label
		WHERE category_id IN 
		(SELECT id FROM category WHERE user_id=$1)
		`
	labels := []Label{}
	err := db.Select(&labels, query, user)
	return &labels, err
}

func (c *Label) Add(user int64) error {
	query := `
		INSERT INTO label 
		(name, description, category_id)
		SELECT
		:name,:description,:category_id
		WHERE :user_id IN (
			SELECT user_id FROM category 
			WHERE id=:category_id)
		`
	type tmp struct {
		Label
		User int64 `db:"user_id"`
	}
	t := &tmp{*c, int64(user)}
	_, err := db.NamedExec(query, t)
	return err
}

func (c *Label) Get(user int64, id int64) (Endpoint, error) {
	query := `
		SELECT id, name, description, category_id
		FROM label
		WHERE id=$2 AND category_id IN (
			SELECT id FROM category 
			WHERE user_id=$1::bigint)
		`

	label := Label{}
	err := db.Get(&label, query, user, id)
	if err != nil {
		return nil, err
	}
	return &label, nil
}

func (c *Label) Del(user int64, id int64) error {
	query := `
		DELETE FROM label WHERE
		id=:id AND category_id IN (
			SELECT id FROM category
			WHERE user_id=:user_id) 
		`
	label := Label{ID: id}
	type tmp struct {
		Label
		User int64 `db:"user_id"`
	}
	t := &tmp{label, user}
	_, err := db.NamedExec(query, t)
	return err
}

func (c *Label) Mod(user int64, id int64) error {
	query := `
		UPDATE label SET
		name = :name,
		description = :description,
		category_id = CASE WHEN :category_id IN (
			SELECT id FROM category
			WHERE user_id=:user_id)
			THEN :category_id ELSE category_id END
		WHERE id=:id
		AND category_id IN (
			SELECT id FROM category
			WHERE user_id=:user_id)
		`
	type tmp struct {
		Label
		User int64 `db:"user_id"`
	}
	c.ID = id
	t := tmp{*c, user}
	_, err := db.NamedExec(query, t)
	return err
}
