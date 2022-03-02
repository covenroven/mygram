CREATE TABLE social_medias (
	id SERIAL PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	social_media_url text NOT NULL,
	user_id int REFERENCES users(id) ON DELETE CASCADE,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL
);
