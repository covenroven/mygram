package api

import (
	"fmt"
	"net/http"

	"github.com/covenroven/mygram/repository"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// Main api application that contains database struct
type App struct {
	db    *sqlx.DB
	repos *repository.Repository
}

// Return internal server error response and print the error
func throwServerError(c *gin.Context, err error) {
	fmt.Println("[Error]", err)
	res := gin.H{
		"errors": "Internal server error",
	}
	c.JSON(http.StatusInternalServerError, res)
}

func NewApp(db *sqlx.DB, repos *repository.Repository) *App {
	return &App{
		db:    db,
		repos: repos,
	}
}

func currentUserID(c *gin.Context) (uint64, bool) {
	var userID uint64
	id, ok := c.Get("userID")
	if !ok {
		return userID, false
	}
	userID = id.(uint64)

	return userID, true
}
