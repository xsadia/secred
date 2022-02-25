package repository

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email        string `json:"email"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token"`
}

func hashPassword(password []byte, salt int) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(password, salt)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func checkEMailAvailability(email string, db *sql.DB) (err error) {
	return db.QueryRow("SELECT COUNT(*) FROM users WHERE users.email = $1", email).Scan(&err)
}

// func (u *User) Create(db *sql.DB) error {

// }
