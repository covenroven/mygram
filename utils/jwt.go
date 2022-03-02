package utils

import (
	"errors"
	"log"

	"github.com/covenroven/mygram/config"
	"github.com/dgrijalva/jwt-go"
)

var (
	errInvalidToken = errors.New("Invalid token")
)

// VerifyJWT checks whether the given token is valid
func VerifyJWT(token string) (interface{}, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errInvalidToken
		}
		return []byte(config.JWT_SECRET_KEY), nil
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Printf("%+v", t.Claims)
	if _, ok := t.Claims.(jwt.MapClaims); !ok && !t.Valid {
		return nil, errInvalidToken
	}

	return t.Claims.(jwt.MapClaims), nil
}

// GenerateJWT returns new token string
func GenerateJWT(id uint64, email string) string {
	claims := jwt.MapClaims{
		"id":    id,
		"email": email,
	}

	parseToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, _ := parseToken.SignedString([]byte(config.JWT_SECRET_KEY))

	return signedToken
}
