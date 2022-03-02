package repository

import (
	"time"

	"github.com/covenroven/mygram/model"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type CommentRepository struct {
	db *sqlx.DB
}

func (cr *CommentRepository) GetAll() ([]model.Comment, error) {
	var comments []model.Comment
	var userIDs []uint64
	var photoIDs []uint64
	query := "SELECT * FROM comments ORDER BY id ASC;"
	rows, err := cr.db.Queryx(query)
	for rows.Next() {
		var comment model.Comment
		err = rows.StructScan(&comment)

		comments = append(comments, comment)
		userIDs = append(userIDs, comment.UserID)
		photoIDs = append(photoIDs, comment.PhotoID)
	}
	if err != nil {
		return comments, err
	}
	defer rows.Close()

	if len(comments) > 0 {
		// Fetch users
		query = `
			SELECT id, email, username
			FROM users
			WHERE id = ANY($1);
		`
		rows, err = cr.db.Queryx(query, pq.Array(userIDs))
		if err != nil {
			return comments, nil
		}
		defer rows.Close()

		users := map[uint64]model.User{}
		for rows.Next() {
			var user model.User
			if err := rows.StructScan(&user); err != nil {
				return comments, err
			}

			users[user.ID] = user
		}

		// Fetch photos
		query = `
			SELECT *
			FROM photos
			WHERE id = ANY($1);
		`
		rows, err = cr.db.Queryx(query, pq.Array(photoIDs))
		if err != nil {
			return comments, nil
		}
		defer rows.Close()

		photos := map[uint64]model.Photo{}
		for rows.Next() {
			var photo model.Photo
			if err := rows.StructScan(&photo); err != nil {
				return comments, err
			}

			photos[photo.ID] = photo
		}

		for i, s := range comments {
			comments[i].User = users[s.UserID]
			comments[i].Photo = photos[s.PhotoID]
		}
	}

	return comments, nil
}

// Create will insert a new data
func (cr *CommentRepository) Create(req model.CommentCreationRequest) (model.Comment, error) {
	var comment model.Comment
	query := `
		INSERT INTO comments (message, photo_id, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING *;
	`
	now := time.Now()
	err := cr.db.QueryRowx(query, req.Message, req.PhotoID, req.UserID, now, now).
		StructScan(&comment)
	if err != nil {
		return comment, err
	}

	return comment, nil
}

// Find will fetch data based on its ID
func (cr *CommentRepository) Find(id uint64) (model.Comment, error) {
	var comment model.Comment
	query := "SELECT * FROM comments WHERE id = $1 LIMIT 1;"
	if err := cr.db.Get(&comment, query, id); err != nil {
		return comment, err
	}

	return comment, nil
}

// Update will update data based on its ID with request data supplied
func (cr *CommentRepository) Update(id uint64, req model.CommentUpdateRequest) (model.Comment, error) {
	var comment model.Comment
	now := time.Now()
	query := `
		UPDATE comments SET message = $1, updated_at = $2
		WHERE id = $3
		RETURNING *;
	`
	err := cr.db.QueryRowx(query, req.Message, now, id).StructScan(&comment)
	if err != nil {
		return comment, err
	}

	return comment, nil
}

// Delete will delete data based on its ID
func (cr *CommentRepository) Delete(id uint64) error {
	query := "DELETE FROM comments WHERE id = $1;"
	_, err := cr.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
