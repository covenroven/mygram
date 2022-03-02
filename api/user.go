package api

import (
	"log"
	"net/http"

	"github.com/covenroven/mygram/model"
	"github.com/covenroven/mygram/utils"
	"github.com/gin-gonic/gin"
)

// RegisterUser will create a new user
func (a *App) RegisterUser(c *gin.Context) {
	var userReq model.UserCreationRequest

	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}
	if err := userReq.UniqueValidation(a.db); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	user, err := a.repos.User.Create(userReq)
	if err != nil {
		log.Println("[RegisterUser]", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"age":      user.Age,
		"email":    user.Email,
		"id":       user.ID,
		"username": user.Username,
	})
}

// LoginUser will check the given email and password, and return JWT if successful
func (a *App) LoginUser(c *gin.Context) {
	var loginReq model.UserLoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	user, err := a.repos.User.FindByEmail(loginReq.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid email or password",
		})
		return
	}
	if !user.VerifyPassword(loginReq.Password) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid email or password",
		})
		return
	}

	token := utils.GenerateJWT(user.ID, user.Email)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// UpdateUser will update user data based on User ID
func (a *App) UpdateUser(c *gin.Context) {
	userID, _ := c.Get("userID")
	id := userID.(uint64)

	user, err := a.repos.User.Find(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "User not found",
		})
		return
	}

	var userReq model.UserUpdateRequest
	userReq.SetUser(user)

	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}
	if err := userReq.UniqueValidation(a.db); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	user, err = a.repos.User.Update(id, userReq)
	if err != nil {
		log.Println("[UpdateUser]", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"username":   user.Username,
		"age":        user.Age,
		"updated_at": user.UpdatedAt,
	})
}

// DeleteCurrentUser will delete user data based on currently logged in data
func (a *App) DeleteCurrentUser(c *gin.Context) {
	userID, _ := c.Get("userID")
	id := userID.(uint64)
	err := a.repos.User.Delete(id)
	if err != nil {
		log.Println("[DeleteCurrentUser]", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your account has been successfully deleted",
	})
}
