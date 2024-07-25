package util

import "golang.org/x/crypto/bcrypt"

func IsCorrectPassword(password string, hashedPassword string) bool {

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false
	}
	return true
}
