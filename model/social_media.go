package model

import "time"

type SocialMedia struct {
	ID             uint64    `json:"id"`
	Name           string    `json:"name"`
	SocialMediaURL string    `json:"social_media_url" db:"social_media_url"`
	UserID         uint64    `json:"user_id" db:"user_id"`
	User           User      `json:"user"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type SocialMediaCreationRequest struct {
	Name           string `json:"name" binding:"required"`
	SocialMediaURL string `json:"social_media_url" binding:"required""`
	UserID         uint64 `json:"user_id"`
}

type SocialMediaUpdateRequest struct {
	Name           string `json:"name"`
	SocialMediaURL string `json:"social_media_url"`
}
