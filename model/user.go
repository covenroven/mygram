package model

import (
	"errors"
	"time"

	"github.com/covenroven/mygram/utils"
	"github.com/jmoiron/sqlx"
)

// Model representation of data from database
type User struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// VerifyPassword checks whether the password given matches with the hashed one
func (u User) VerifyPassword(password string) bool {
	return utils.CheckPasswordHash(u.Password, password)
}

// Struct for user registration request
type UserCreationRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Age      int    `json:"age" binding:"required,min=9"`
}

// UniqueValidation makes sure the supplied fields is unique in the database
func (ucr UserCreationRequest) UniqueValidation(db *sqlx.DB) error {
	if ucr.Username != "" {
		if utils.DBValidateExists(db, ucr.Username, "users", "username") {
			return errors.New("Username already exists")
		}
	}

	if ucr.Email != "" {
		if utils.DBValidateExists(db, ucr.Email, "users", "email") {
			return errors.New("Email already exists")
		}
	}

	return nil
}

// HashedPassword returns bcrypt-hashed password in string
func (ucr UserCreationRequest) HashedPassword() string {
	return utils.HashPassword(ucr.Password)
}

// Struct for user update request
type UserUpdateRequest struct {
	user     User
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// UniqueValidation makes sure the supplied fields is unique in the database
func (uur UserUpdateRequest) UniqueValidation(db *sqlx.DB) error {
	if (uur.Username != "") && (uur.Username != uur.user.Username) {
		if utils.DBValidateExists(db, uur.Username, "users", "username") {
			return errors.New("Username already exists")
		}
	}

	if (uur.Username != "") && (uur.Email != uur.user.Email) {
		if utils.DBValidateExists(db, uur.Email, "users", "email") {
			return errors.New("Email already exists")
		}
	}

	return nil
}

// SetUser sets user that is being updated to the struct
func (uur *UserUpdateRequest) SetUser(user User) {
	uur.user = user
}

// Struct for user login request
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
