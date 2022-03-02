package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/covenroven/mygram/model"
	"github.com/gin-gonic/gin"
)

// CreateSocialMedia creates new socialMedia for the current user
func (a *App) CreateSocialMedia(c *gin.Context) {
	var socialMediaReq model.SocialMediaCreationRequest

	if err := c.ShouldBindJSON(&socialMediaReq); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	userID, _ := currentUserID(c)
	socialMediaReq.UserID = userID

	socialMedia, err := a.repos.SocialMedia.Create(socialMediaReq)
	if err != nil {
		log.Println("[CreateSocialMedia]", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":               socialMedia.ID,
		"name":             socialMedia.Name,
		"social_media_url": socialMedia.SocialMediaURL,
		"user_id":          socialMedia.UserID,
		"created_at":       socialMedia.CreatedAt,
	})
}

// GetSocialMedias fetches all socialMedias of the current user
func (a *App) GetSocialMedias(c *gin.Context) {
	socialMedias, err := a.repos.SocialMedia.GetAll()
	if err != nil {
		log.Println("[GetSocialMedia]", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	resp := []map[string]interface{}{}
	for _, p := range socialMedias {
		resp = append(resp, map[string]interface{}{
			"id":               p.ID,
			"name":             p.Name,
			"social_media_url": p.SocialMediaURL,
			"user_id":          p.UserID,
			"created_at":       p.CreatedAt,
			"updated_at":       p.UpdatedAt,
			"User": map[string]interface{}{
				"email": p.User.Email,
				"id":    p.User.ID,
			},
		})
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateSocialMedia updates a current user's socialMedia based on its ID
func (a *App) UpdateSocialMedia(c *gin.Context) {
	userID, _ := currentUserID(c)

	id, err := strconv.ParseUint(c.Param("socialMediaID"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Social media not found",
		})
		return
	}

	socialMedia, err := a.repos.SocialMedia.Find(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Social media not found",
		})
		return
	}

	if socialMedia.UserID != userID {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Social media not found",
		})
		return
	}

	var socialMediaReq model.SocialMediaUpdateRequest

	if err := c.ShouldBindJSON(&socialMediaReq); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	socialMedia, err = a.repos.SocialMedia.Update(id, socialMediaReq)
	if err != nil {
		log.Println("[UpdateSocialMedia]", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":               socialMedia.ID,
		"name":             socialMedia.Name,
		"social_media_url": socialMedia.SocialMediaURL,
		"user_id":          socialMedia.UserID,
		"updated_at":       socialMedia.UpdatedAt,
	})
}

// DeleteSocialMedia deletes a current user's socialMedia based on its ID
func (a *App) DeleteSocialMedia(c *gin.Context) {
	userID, _ := currentUserID(c)

	id, err := strconv.ParseUint(c.Param("socialMediaID"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Social media not found",
		})
		return
	}

	socialMedia, err := a.repos.SocialMedia.Find(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Social media not found",
		})
		return
	}

	if socialMedia.UserID != userID {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Social media not found",
		})
		return
	}

	err = a.repos.SocialMedia.Delete(id)
	if err != nil {
		log.Println("[DeleteSocialMedia]", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your social media has been successfully deleted",
	})
}
