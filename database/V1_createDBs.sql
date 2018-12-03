DROP DATABASE IF EXISTS pefi;
CREATE DATABASE pefi;

\c pefi

CREATE TABLE IF NOT EXISTS accounts (
	id uuid PRIMARY KEY,
	name VARCHAR(40),
	description VARCHAR(255),
	owner_id uuid,
	balance INTEGER,
	currency VARCHAR(3)
);

CREATE TABLE IF NOT EXISTS users (
	id uuid PRIMARY KEY,
	name VARCHAR(40)
);
