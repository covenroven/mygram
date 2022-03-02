CREATE TABLE photos (
	id SERIAL PRIMARY KEY,
	title varchar(255) NOT NULL,
	caption text NOT NULL,
	photo_url text NOT NULL,
	user_id int REFERENCES users(id) ON DELETE CASCADE,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL
);
