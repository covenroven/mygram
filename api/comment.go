package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/covenroven/mygram/model"
	"github.com/gin-gonic/gin"
)

// CreateComment creates new comment for the current user
func (a *App) CreateComment(c *gin.Context) {
	var commentReq model.CommentCreationRequest

	if err := c.ShouldBindJSON(&commentReq); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	userID, _ := currentUserID(c)
	commentReq.UserID = userID

	if err := commentReq.ExistValidation(a.db); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	comment, err := a.repos.Comment.Create(commentReq)
	if err != nil {
		log.Println("[CreateComment]", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         comment.ID,
		"message":    comment.Message,
		"photo_id":   comment.PhotoID,
		"user_id":    comment.UserID,
		"created_at": comment.CreatedAt,
	})
}

// GetComments fetches all comments of the current user
func (a *App) GetComments(c *gin.Context) {
	comments, err := a.repos.Comment.GetAll()
	if err != nil {
		log.Println("[GetComments]", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	resp := []map[string]interface{}{}
	for _, p := range comments {
		data := map[string]interface{}{
			"id":         p.ID,
			"message":    p.Message,
			"photo_id":   p.PhotoID,
			"user_id":    p.UserID,
			"created_at": p.CreatedAt,
			"updated_at": p.UpdatedAt,
			"User": map[string]interface{}{
				"username": p.User.Username,
				"email":    p.User.Email,
				"id":       p.User.ID,
			},
			"Photo": map[string]interface{}{
				"id":        p.Photo.ID,
				"title":     p.Photo.Title,
				"caption":   p.Photo.Caption,
				"photo_url": p.Photo.PhotoURL,
				"user_id":   p.Photo.UserID,
			},
		}
		resp = append(resp, data)
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateComment updates a current user's comment based on its ID
func (a *App) UpdateComment(c *gin.Context) {
	userID, _ := currentUserID(c)

	id, err := strconv.ParseUint(c.Param("commentID"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Comment not found",
		})
		return
	}

	comment, err := a.repos.Comment.Find(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Comment not found",
		})
		return
	}

	if comment.UserID != userID {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Comment not found",
		})
		return
	}

	var commentReq model.CommentUpdateRequest

	if err := c.ShouldBindJSON(&commentReq); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	comment, err = a.repos.Comment.Update(id, commentReq)
	if err != nil {
		log.Println("[UpdateComment]", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         comment.ID,
		"message":    comment.Message,
		"photo_id":   comment.PhotoID,
		"user_id":    comment.UserID,
		"updated_at": comment.UpdatedAt,
	})
}

// DeleteComment deletes a current user's comment based on its ID
func (a *App) DeleteComment(c *gin.Context) {
	userID, _ := currentUserID(c)

	id, err := strconv.ParseUint(c.Param("commentID"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Comment not found",
		})
		return
	}

	comment, err := a.repos.Comment.Find(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Comment not found",
		})
		return
	}

	if comment.UserID != userID {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Comment not found",
		})
		return
	}

	err = a.repos.Comment.Delete(id)
	if err != nil {
		log.Println("[DeleteComment]", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your comment has been successfully deleted",
	})
}
