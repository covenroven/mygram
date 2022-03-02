package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/covenroven/mygram/model"
	"github.com/gin-gonic/gin"
)

// CreatePhoto creates new photo for the current user
func (a *App) CreatePhoto(c *gin.Context) {
	var photoReq model.PhotoCreationRequest

	if err := c.ShouldBindJSON(&photoReq); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	userID, _ := currentUserID(c)
	photoReq.UserID = userID

	photo, err := a.repos.Photo.Create(photoReq)
	if err != nil {
		log.Println("[CreatePhoto]", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         photo.ID,
		"title":      photo.Title,
		"caption":    photo.Caption,
		"photo_url":  photo.PhotoURL,
		"user_id":    photo.UserID,
		"created_at": photo.CreatedAt,
	})
}

// GetPhotos fetches all photos of the current user
func (a *App) GetPhotos(c *gin.Context) {
	photos, err := a.repos.Photo.GetAll()
	if err != nil {
		log.Println("[GetPhotos]", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	var resp []map[string]interface{}
	for _, p := range photos {
		resp = append(resp, map[string]interface{}{
			"id":         p.ID,
			"title":      p.Title,
			"caption":    p.Caption,
			"photo_url":  p.PhotoURL,
			"user_id":    p.UserID,
			"created_at": p.CreatedAt,
			"updated_at": p.UpdatedAt,
			"User": map[string]interface{}{
				"email":    p.User.Email,
				"username": p.User.Username,
			},
		})
	}

	c.JSON(http.StatusOK, resp)
}

// UpdatePhoto updates a current user's photo based on its ID
func (a *App) UpdatePhoto(c *gin.Context) {
	userID, _ := currentUserID(c)

	id, err := strconv.ParseUint(c.Param("photoID"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Photo not found",
		})
		return
	}

	photo, err := a.repos.Photo.Find(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Photo not found",
		})
		return
	}

	if photo.UserID != userID {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Photo not found",
		})
		return
	}

	var photoReq model.PhotoUpdateRequest

	if err := c.ShouldBindJSON(&photoReq); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	photo, err = a.repos.Photo.Update(id, photoReq)
	if err != nil {
		log.Println("[UpdatePhotos]", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         photo.ID,
		"title":      photo.Title,
		"caption":    photo.Caption,
		"photo_url":  photo.PhotoURL,
		"user_id":    photo.UserID,
		"updated_at": photo.UpdatedAt,
	})
}

// DeletePhoto deletes a current user's photo based on its ID
func (a *App) DeletePhoto(c *gin.Context) {
	userID, _ := currentUserID(c)

	id, err := strconv.ParseUint(c.Param("photoID"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Photo not found",
		})
		return
	}

	photo, err := a.repos.Photo.Find(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Photo not found",
		})
		return
	}

	if photo.UserID != userID {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Photo not found",
		})
		return
	}

	err = a.repos.Photo.Delete(id)
	if err != nil {
		log.Println("[DeletePhoto]", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your photo has been successfully deleted",
	})
}
