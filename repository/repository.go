package repository

import "github.com/jmoiron/sqlx"

// Repository collects every repositories in single struct
type Repository struct {
	User        *UserRepository
	Photo       *PhotoRepository
	Comment     *CommentRepository
	SocialMedia *SocialMediaRepository
}

// NewRepository creates every repository and return the collection of them
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User: &UserRepository{
			db: db,
		},
		Photo: &PhotoRepository{
			db: db,
		},
		Comment: &CommentRepository{
			db: db,
		},
		SocialMedia: &SocialMediaRepository{
			db: db,
		},
	}
}
