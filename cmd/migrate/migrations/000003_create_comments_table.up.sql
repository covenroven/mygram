CREATE TABLE comments (
	id SERIAL PRIMARY KEY,
	user_id int REFERENCES users(id) ON DELETE CASCADE,
	photo_id int DEFAULT NULL REFERENCES photos(id) ON DELETE CASCADE,
	message text NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL
);
