package internal

import "golang.org/x/crypto/bcrypt"

func HashPassword(password []byte, salt int) string {
	hash, _ := bcrypt.GenerateFromPassword(password, salt)

	return string(hash)
}
