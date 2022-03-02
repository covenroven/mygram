CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	username varchar(255) NOT NULL UNIQUE,
	email varchar(255) NOT NULL UNIQUE,
	password varchar(255) NOT NULL,
	age int NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL
);
