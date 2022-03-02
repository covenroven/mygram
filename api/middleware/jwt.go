package middleware

import (
	"errors"
	"strings"

	"github.com/covenroven/mygram/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	errUnauthorized = errors.New("Unauthorized")
	errInvalidToken = errors.New("Invalid token")
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if ok := strings.HasPrefix(authHeader, "Bearer"); !ok {
			c.AbortWithStatusJSON(401, gin.H{
				"message": errUnauthorized.Error(),
			})
			return
		}

		stringToken := strings.Split(authHeader, " ")[1]
		token, err := utils.VerifyJWT(stringToken)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.Set("userID", uint64(token.(jwt.MapClaims)["id"].(float64)))
		c.Set("email", token.(jwt.MapClaims)["email"].(string))

		c.Next()
	}
}
