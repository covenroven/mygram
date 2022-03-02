package model

import (
	"errors"
	"time"

	"github.com/covenroven/mygram/utils"
	"github.com/jmoiron/sqlx"
)

type Comment struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id" db:"user_id"`
	PhotoID   uint64    `json:"photo_id" db:"photo_id"`
	Message   string    `json:"message"`
	User      User      `json:"user"`
	Photo     Photo     `json:"photo"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CommentCreationRequest struct {
	Message string `json:"message" binding:"required"`
	UserID  uint64 `json:"user_id"`
	PhotoID uint64 `json:"photo_id" binding:"required"`
}

func (ccr CommentCreationRequest) ExistValidation(db *sqlx.DB) error {
	if ccr.PhotoID > 0 {
		if !utils.DBValidateExists(db, ccr.PhotoID, "photos", "id") {
			return errors.New("Photo doesn't exist")
		}
	}

	return nil
}

type CommentUpdateRequest struct {
	Message string `json:"message"`
}
