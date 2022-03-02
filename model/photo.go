package model

import "time"

type Photo struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	PhotoURL  string    `json:"photo_url" db:"photo_url"`
	UserID    uint64    `json:"user_id" db:"user_id"`
	User      User      `json:"user"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type PhotoCreationRequest struct {
	UserID   uint64
	Title    string `json:"title" binding:"required"`
	Caption  string `json:"caption"`
	PhotoURL string `json:"photo_url" binding:"required"`
}

type PhotoUpdateRequest struct {
	UserID   uint64
	Title    string `json:"title" binding:"required"`
	Caption  string `json:"caption"`
	PhotoURL string `json:"photo_url" binding:"required"`
}
