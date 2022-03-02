package repository

import (
	"time"

	"github.com/covenroven/mygram/model"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

// Create will insert a new data
func (ur *UserRepository) Create(req model.UserCreationRequest) (model.User, error) {
	var user model.User
	query := `
		INSERT INTO users (username, email, password, age, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING *;
	`
	now := time.Now()
	err := ur.db.QueryRowx(query, req.Username, req.Email, req.HashedPassword(), req.Age, now, now).
		StructScan(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

// Find will fetch data based on its ID
func (ur *UserRepository) Find(id uint64) (model.User, error) {
	var user model.User
	query := "SELECT * FROM users WHERE id = $1 LIMIT 1;"
	if err := ur.db.Get(&user, query, id); err != nil {
		return user, err
	}

	return user, nil
}

// Update will update data based on its ID with request data supplied
func (ur *UserRepository) Update(id uint64, req model.UserUpdateRequest) (model.User, error) {
	var user model.User
	now := time.Now()
	query := `
		UPDATE users SET username = $1, email = $2, updated_at = $3
		WHERE id = $4
		RETURNING *;
	`
	err := ur.db.QueryRowx(query, req.Username, req.Email, now, id).StructScan(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

// Delete will delete data based on its ID
func (ur *UserRepository) Delete(id uint64) error {
	query := "DELETE FROM users WHERE id = $1;"
	_, err := ur.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

// FindByEmail will return data via email lookup
func (ur *UserRepository) FindByEmail(email string) (model.User, error) {
	var user model.User
	query := "SELECT * FROM users WHERE email = $1 LIMIT 1;"
	if err := ur.db.QueryRowx(query, email).StructScan(&user); err != nil {
		return user, err
	}

	return user, nil
}
