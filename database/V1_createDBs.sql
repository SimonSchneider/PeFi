DROP DATABASE IF EXISTS pefi;
CREATE DATABASE pefi;

\c pefi

CREATE TABLE IF NOT EXISTS users (
	id BIGSERIAL NOT NULL PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS category (
	id BIGSERIAL NOT NULL PRIMARY KEY,
	name TEXT,
	description TEXT,
	user_id BIGINT NOT NULL REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS external_account (
	id BIGSERIAL NOT NULL PRIMARY KEY,
	name VARCHAR(40),
	description VARCHAR(255),
	category_id BIGINT NOT NULL REFERENCES category(id)
);

CREATE TABLE IF NOT EXISTS internal_account (
	external_account_id BIGINT NOT NULL REFERENCES external_account(id) ON DELETE CASCADE,
	balance DECIMAL (18, 2)
);

CREATE TABLE IF NOT EXISTS label (
	id BIGSERIAL NOT NULL PRIMARY KEY,
	name VARCHAR(40),
	description VARCHAR(255),
	category_id BIGINT NOT NULL REFERENCES category(id)
);

CREATE TABLE IF NOT EXISTS transaction (
	id BIGSERIAL NOT NULL PRIMARY KEY,
	time TIMESTAMP NOT NULL,
	amount DECIMAL NOT NULL,
	sender_id BIGINT NOT NULL REFERENCES external_account(id),
	receiver_id BIGINT NOT NULL REFERENCES external_account(id),
	label_id BIGINT NOT NULL REFERENCES label(id)
);

CREATE TABLE IF NOT EXISTS loan (
	transaction_id BIGINT NOT NULL REFERENCES transaction(id) ON DELETE CASCADE,
	payback_id BIGINT NOT NULL REFERENCES transaction(id),
	PRIMARY KEY (transaction_id, payback_id)
);
