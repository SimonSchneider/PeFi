package models

type (
	ExternalAccount struct {
		ID          int64  `json:"id" db:"id"`
		Name        string `json:"name" db:"name"`
		Description string `json:"description" db:"description"`
		CategoryID  int64  `json:"category_id" db:"category_id"`
	}
)

// Name returns the name of category endpoint
func (c *ExternalAccount) URL() string {
	return "/accounts/external"
}

func (c *ExternalAccount) GetAll(user int64) (interface{}, error) {
	query := `
		SELECT id, name, description, category_id
		FROM external_account
		WHERE category_id IN 
		(SELECT id FROM category WHERE user_id=$1)
		AND id NOT IN (SELECT external_account_id
			FROM internal_account)
		`
	categories := []ExternalAccount{}
	err := db.Select(&categories, query, user)
	return &categories, err
}

func (c *ExternalAccount) Add(user int64) error {
	query := `
		INSERT INTO external_account 
		(name, description, category_id)
		SELECT
		:name,:description,:category_id
		WHERE :user_id IN (
			SELECT user_id FROM category 
			WHERE id=:category_id)
		`
	type tmp struct {
		ExternalAccount
		User int64 `db:"user_id"`
	}
	t := &tmp{*c, int64(user)}
	_, err := db.NamedExec(query, t)
	return err
}

func (c *ExternalAccount) Get(user int64, id int64) (Endpoint, error) {
	query := `
		SELECT id, name, description, category_id
		FROM external_account
		WHERE id=$2 AND category_id IN (
			SELECT id FROM category 
			WHERE user_id=$1::bigint)
			AND id NOT IN (SELECT external_account_id
			FROM internal_account)
		`

	category := ExternalAccount{}
	err := db.Get(&category, query, user, id)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (c *ExternalAccount) Del(user int64, id int64) error {
	query := `
		DELETE FROM external_account WHERE
		id=:id AND category_id IN (
			SELECT id FROM category
			WHERE user_id=:user_id) 
			AND id NOT IN (SELECT external_account_id
			FROM internal_account)
		`
	category := ExternalAccount{ID: id}
	type tmp struct {
		ExternalAccount
		User int64 `db:"user_id"`
	}
	t := &tmp{category, user}
	_, err := db.NamedExec(query, t)
	return err
}

func (c *ExternalAccount) Mod(user int64, id int64) error {
	query := `
		UPDATE external_account SET
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
			AND id NOT IN (SELECT external_account_id
			FROM internal_account)
		`
	type tmp struct {
		ExternalAccount
		User int64 `db:"user_id"`
	}
	c.ID = id
	t := tmp{*c, user}
	_, err := db.NamedExec(query, t)
	return err
}
