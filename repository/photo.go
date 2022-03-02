package repository

import (
	"time"

	"github.com/covenroven/mygram/model"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type PhotoRepository struct {
	db *sqlx.DB
}

func (pr *PhotoRepository) GetAll() ([]model.Photo, error) {
	var photos []model.Photo
	var userIDs []uint64
	query := "SELECT * FROM photos ORDER BY id ASC;"
	rows, err := pr.db.Queryx(query)
	for rows.Next() {
		var photo model.Photo
		err = rows.StructScan(&photo)

		photos = append(photos, photo)
		userIDs = append(userIDs, photo.UserID)
	}
	if err != nil {
		return photos, err
	}
	defer rows.Close()

	// Fetch users
	if len(photos) > 0 {
		query = `
			SELECT id, email, username
			FROM users
			WHERE id = ANY($1);
		`
		rows, err = pr.db.Queryx(query, pq.Array(userIDs))
		if err != nil {
			return photos, nil
		}
		defer rows.Close()

		users := map[uint64]model.User{}
		for rows.Next() {
			var user model.User
			if err := rows.StructScan(&user); err != nil {
				return photos, err
			}

			users[user.ID] = user
		}

		for i, p := range photos {
			photos[i].User = users[p.UserID]
		}
	}

	return photos, nil
}

// Create will insert a new data
func (pr *PhotoRepository) Create(req model.PhotoCreationRequest) (model.Photo, error) {
	var photo model.Photo
	query := `
		INSERT INTO photos (title, caption, photo_url, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING *;
	`
	now := time.Now()
	err := pr.db.QueryRowx(query, req.Title, req.Caption, req.PhotoURL, req.UserID, now, now).
		StructScan(&photo)
	if err != nil {
		return photo, err
	}

	return photo, nil
}

// Find will fetch data based on its ID
func (pr *PhotoRepository) Find(id uint64) (model.Photo, error) {
	var photo model.Photo
	query := "SELECT * FROM photos WHERE id = $1 LIMIT 1;"
	if err := pr.db.Get(&photo, query, id); err != nil {
		return photo, err
	}

	return photo, nil
}

// Update will update data based on its ID with request data supplied
func (pr *PhotoRepository) Update(id uint64, req model.PhotoUpdateRequest) (model.Photo, error) {
	var photo model.Photo
	now := time.Now()
	query := `
		UPDATE photos SET title = $1, caption = $2, photo_url = $3, updated_at = $4
		WHERE id = $5
		RETURNING *;
	`
	err := pr.db.QueryRowx(query, req.Title, req.Caption, req.PhotoURL, now, id).StructScan(&photo)
	if err != nil {
		return photo, err
	}

	return photo, nil
}

// Delete will delete data based on its ID
func (pr *PhotoRepository) Delete(id uint64) error {
	query := "DELETE FROM photos WHERE id = $1;"
	_, err := pr.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
