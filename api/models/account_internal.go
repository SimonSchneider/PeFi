package models

import "fmt"

type (
	InternalAccount struct {
		ExternalAccount
		Balance float64 `json:"balance" db:"balance"`
	}
)

// Name returns the name of category endpoint
func (c *InternalAccount) URL() string {
	return "/accounts/internal"
}

func (c *InternalAccount) GetAll(user int64) (interface{}, error) {
	query := `
		SELECT id, name, description, category_id, balance
		FROM external_account, internal_account
		WHERE id = external_account_id AND category_id IN
		(SELECT id FROM category WHERE user_id=$1)
		`
	categories := []InternalAccount{}
	err := db.Select(&categories, query, user)
	fmt.Println(categories)
	return &categories, err
}

func (c *InternalAccount) Add(user int64) error {
	query := `
	WITH ins1 AS ( INSERT INTO 
		external_account(name, description, category_id)
		SELECT :name, :description, :category_id 
		WHERE :category_id IN (SELECT id 
			FROM category WHERE user_id=:user_id)
		RETURNING id)
	INSERT INTO internal_account(external_account_id, balance)
	VALUES ( (SELECT id from ins1), :balance);  
	`
	type tmp struct {
		InternalAccount
		User int64 `db:"user_id"`
	}
	t := &tmp{*c, int64(user)}
	_, err := db.NamedExec(query, t)
	return err
}

func (c *InternalAccount) Get(user int64, id int64) (Endpoint, error) {
	query := `
		SELECT id, name, description, category_id, balance
		FROM external_account, internal_account
		WHERE id = external_account_id AND id=$2
		AND category_id IN 
			(SELECT id FROM category WHERE user_id=$1)
		`

	category := InternalAccount{}
	err := db.Get(&category, query, user, id)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (c *InternalAccount) Del(user int64, id int64) error {
	query := `
		DELETE FROM internal_account WHERE
		external_account_id IN (
			SELECT id FROM external_account
			WHERE id=:id AND category_id IN
			(
				SELECT id FROM category
				WHERE user_id=:user_id
			)
		)
		`
	acc := InternalAccount{ExternalAccount: ExternalAccount{ID: id}}
	type tmp struct {
		InternalAccount
		User int64 `db:"user_id"`
	}
	t := &tmp{acc, user}
	_, err := db.NamedExec(query, t)
	return err
}

func (c *InternalAccount) Mod(user int64, id int64) error {
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
			AND id IN (SELECT external_account_id
			FROM internal_account)
		`
	type tmp struct {
		InternalAccount
		User int64 `db:"user_id"`
	}
	c.ID = id
	t := tmp{*c, user}
	_, err := db.NamedExec(query, t)
	if err != nil {
		return err
	}
	query = `
	UPDATE internal_account SET balance = :balance
	WHERE external_account_id=:id
	AND external_account_id in (
		SELECT id FROM
		external_account WHERE
		category_id IN (
			SELECT id FROM category
			WHERE user_id=:user_id
		)
	)
	`
	_, err = db.NamedExec(query, t)
	return err
}
