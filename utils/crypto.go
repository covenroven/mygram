package utils

import (
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 10

// HashPassword returns hashed password of the given string
func HashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)

	return string(hash)
}

// CheckPasswordHash checks whether the password matchs the given hashed string
func CheckPasswordHash(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}

	return true
}
