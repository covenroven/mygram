package repository

import (
	"time"

	"github.com/covenroven/mygram/model"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type SocialMediaRepository struct {
	db *sqlx.DB
}

func (smr *SocialMediaRepository) GetAll() ([]model.SocialMedia, error) {
	var socialMedias []model.SocialMedia
	var userIDs []uint64
	query := "SELECT * FROM social_medias ORDER BY id ASC;"
	rows, err := smr.db.Queryx(query)
	for rows.Next() {
		var socialMedia model.SocialMedia
		err = rows.StructScan(&socialMedia)

		socialMedias = append(socialMedias, socialMedia)
		userIDs = append(userIDs, socialMedia.UserID)
	}
	if err != nil {
		return socialMedias, err
	}
	defer rows.Close()

	// Fetch users
	if len(socialMedias) > 0 {
		query = `
			SELECT id, email, username
			FROM users
			WHERE id = ANY($1);
		`
		rows, err = smr.db.Queryx(query, pq.Array(userIDs))
		if err != nil {
			return socialMedias, nil
		}
		defer rows.Close()

		users := map[uint64]model.User{}
		for rows.Next() {
			var user model.User
			if err := rows.StructScan(&user); err != nil {
				return socialMedias, err
			}

			users[user.ID] = user
		}

		for i, s := range socialMedias {
			socialMedias[i].User = users[s.UserID]
		}
	}

	return socialMedias, nil
}

// Create will insert a new data
func (smr *SocialMediaRepository) Create(req model.SocialMediaCreationRequest) (model.SocialMedia, error) {
	var socialMedia model.SocialMedia
	query := `
		INSERT INTO social_medias (name, social_media_url, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING *;
	`
	now := time.Now()
	err := smr.db.QueryRowx(query, req.Name, req.SocialMediaURL, req.UserID, now, now).
		StructScan(&socialMedia)
	if err != nil {
		return socialMedia, err
	}

	return socialMedia, nil
}

// Find will fetch data based on its ID
func (smr *SocialMediaRepository) Find(id uint64) (model.SocialMedia, error) {
	var socialMedia model.SocialMedia
	query := "SELECT * FROM social_medias WHERE id = $1 LIMIT 1;"
	if err := smr.db.Get(&socialMedia, query, id); err != nil {
		return socialMedia, err
	}

	return socialMedia, nil
}

// Update will update data based on its ID with request data supplied
func (smr *SocialMediaRepository) Update(id uint64, req model.SocialMediaUpdateRequest) (model.SocialMedia, error) {
	var socialMedia model.SocialMedia
	now := time.Now()
	query := `
		UPDATE social_medias SET name = $1, social_media_url = $2, updated_at = $3
		WHERE id = $4
		RETURNING *;
	`
	err := smr.db.QueryRowx(query, req.Name, req.SocialMediaURL, now, id).StructScan(&socialMedia)
	if err != nil {
		return socialMedia, err
	}

	return socialMedia, nil
}

// Delete will delete data based on its ID
func (smr *SocialMediaRepository) Delete(id uint64) error {
	query := "DELETE FROM social_medias WHERE id = $1;"
	_, err := smr.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
