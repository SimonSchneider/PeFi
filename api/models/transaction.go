package models

import (
	"time"
)

type (
	Transaction struct {
		ID         int64     `json:"id" db:"id"`
		Time       time.Time `json:"time" db:"time"`
		Amount     float64   `json:"amount,number" db:"amount"`
		SenderID   float64   `json:"sender_id" db:"sender_id"`
		ReceiverID float64   `json:"receiver_id" db:"receiver_id"`
		LabelID    int64     `json:"label_id" db:"label_id"`
	}
)

// Name returns the name of category endpoint
func (c *Transaction) URL() string {
	return "/transactions"
}

func (c *Transaction) GetAll(user int64) (interface{}, error) {
	query := `
	SELECT id, time, amount, sender_id, receiver_id, label_id
	FROM transaction
	WHERE label_id IN (
		SELECT id FROM label 
		WHERE category_id IN (
			SELECT id FROM category
			WHERE user_id=$1
		)
	)
	`
	transactions := []Transaction{}
	err := db.Select(&transactions, query, user)
	return &transactions, err
}

func (c *Transaction) Add(user int64) error {
	query1 := `
	INSERT INTO transaction 
	(time, amount, sender_id, receiver_id, label_id)
	SELECT
	:time,:amount,:sender_id,:receiver_id,:label_id
	WHERE :label_id IN (
		SELECT id FROM label 
		WHERE category_id IN (
			SELECT id FROM category
			WHERE user_id=:user_id
		)
	) AND :receiver_id IN (
		SELECT id FROM external_account
		WHERE category_id IN (
			SELECT id FROM category
			WHERE user_id=:user_id
		)
	) AND :sender_id IN (
		SELECT id FROM external_account
		WHERE category_id IN (
			SELECT id FROM category
			WHERE user_id=:user_id
		)
	)
	`
	query2 := `
	UPDATE internal_account SET balance = balance-:amount
	WHERE external_account_id=:sender_id
	AND :label_id IN (
		SELECT id FROM label 
		WHERE category_id IN (
			SELECT id FROM category
			WHERE user_id=:user_id
		)
	) AND :receiver_id IN (
		SELECT id FROM external_account
		WHERE category_id IN (
			SELECT id FROM category
			WHERE user_id=:user_id
		)
	) AND :sender_id IN (
		SELECT id FROM external_account
		WHERE category_id IN (
			SELECT id FROM category
			WHERE user_id=:user_id
		)
	)
	`
	query3 := `
	UPDATE internal_account SET balance = balance+:amount
	WHERE external_account_id=:receiver_id
	AND :label_id IN (
		SELECT id FROM label 
		WHERE category_id IN (
			SELECT id FROM category
			WHERE user_id=:user_id
		)
	) AND :receiver_id IN (
		SELECT id FROM external_account
		WHERE category_id IN (
			SELECT id FROM category
			WHERE user_id=:user_id
		)
	) AND :sender_id IN (
		SELECT id FROM external_account
		WHERE category_id IN (
			SELECT id FROM category
			WHERE user_id=:user_id
		)
	)
	`
	type tmp struct {
		Transaction
		User int64 `db:"user_id"`
	}
	t := &tmp{*c, int64(user)}
	tx, err := db.Beginx()
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.NamedExec(query1, t)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.NamedExec(query2, t)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.NamedExec(query3, t)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

func (c *Transaction) Get(user int64, id int64) (Endpoint, error) {
	query := `
	SELECT id, time, amount, sender_id, receiver_id, label_id
	FROM transaction
	WHERE id=$2 AND label_id IN (
		SELECT id FROM label 
		WHERE category_id IN (
			SELECT id FROM category
			WHERE user_id=$1
		)
	)
	`
	transaction := Transaction{}
	err := db.Get(&transaction, query, user, id)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (c *Transaction) Del(user int64, id int64) error {
	query1 := `
	DELETE FROM transaction
	WHERE id=:id AND label_id IN (
		SELECT id FROM label 
		WHERE category_id IN (
			SELECT id FROM category
			WHERE user_id=:user_id
		)
	)
	`
	query2 := `
	UPDATE internal_account SET balance = balance+:amount
	WHERE external_account_id IN (
		SELECT sender_id FROM transaction
		WHERE id=:id AND label_id IN (
			SELECT id FROM label 
			WHERE category_id IN (
				SELECT id FROM category
				WHERE user_id=:user_id
			)
		)
	)
	`
	query3 := `
	UPDATE internal_account SET balance = balance-:amount
	WHERE external_account_id IN (
		SELECT receiver_id FROM transaction
		WHERE id=:id AND label_id IN (
			SELECT id FROM label 
			WHERE category_id IN (
				SELECT id FROM category
				WHERE user_id=:user_id
			)
		)
	)
	`
	transaction := Transaction{ID: id}
	type tmp struct {
		Transaction
		User int64 `db:"user_id"`
	}
	t := &tmp{transaction, user}
	tx, err := db.Beginx()
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.NamedExec(query3, t)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.NamedExec(query2, t)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.NamedExec(query1, t)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

func (c *Transaction) Mod(user int64, id int64) error {
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
		Transaction
		User int64 `db:"user_id"`
	}
	c.ID = id
	t := tmp{*c, user}
	_, err := db.NamedExec(query, t)
	return err
}
