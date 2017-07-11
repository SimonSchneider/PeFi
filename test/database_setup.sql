DROP DATABASE IF EXISTS pefi_test;
CREATE DATABASE pefi_test;

\c pefi_test

\i /pefi/database/design.sql

INSERT INTO users (id) VALUES (default), (default);

-- SELECT * FROM users;

INSERT INTO category (name, description, user_id) 
VALUES 	('test1', 'test1', 1),
	('test2', 'test2', 2);

-- SELECT * FROM category;

INSERT INTO external_account(name, description, category_id) 
VALUES 	('test1', 'test1', 1),
	('test2', 'test2', 2),
	('test1', 'test1', 1),
	('test2', 'test2', 2);

-- SELECT * FROM external_account;

INSERT INTO internal_account(external_account_id, balance) 
VALUES 	(1, 1.0),
	(2, 2.0);

-- SELECT * FROM internal_account AS i
	-- JOIN external_account AS e
	-- ON i.external_account_id = e.id;

INSERT INTO label(name, description, category_id)
VALUES 	('test1', 'test1', 1),
	('test2', 'test2', 2);

-- SELECT * FROM label;

INSERT INTO transaction(time, amount, sender_id, receiver_id, label_id)
VALUES 	('2017-07-09 22:53:43', 1.0, 1, 3, 1),
	('2017-07-09 22:53:43', 1.0, 3, 1, 1),
	('2017-07-09 22:53:43', 1.0, 2, 4, 2),
	('2017-07-09 22:53:43', 1.0, 4, 2, 2);

-- SELECT * FROM transaction;

INSERT INTO loan(transaction_id, payback_id)
VALUES 	(1, 3),
	(1, 1),
	(2, 2),
	(2, 4);

-- SELECT * FROM loan AS l JOIN transaction AS t on l.transaction_id=t.id;
