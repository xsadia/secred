package repository

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email        string `json:"email"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token"`
}

func hashPassword(password []byte, salt int) string {
	hash, _ := bcrypt.GenerateFromPassword(password, salt)

	return string(hash)
}

func checkEMailAvailability(email string, db *sql.DB) (err error) {
	return db.QueryRow("SELECT * FROM users WHERE users.email = $1", email).Scan(&err)
}

func (u *User) Create(db *sql.DB) error {
	err := checkEMailAvailability(u.Email, db)

	if err != nil && err != sql.ErrNoRows {
		return errors.New("e-mail already in use")
	}

	u.Password = hashPassword([]byte(u.Password), 8)

	_, err =
		db.Exec(
			"INSERT INTO users (email, username, password) VALUES ($1, $2, $3)",
			u.Email, u.Username, u.Password,
		)

	return err
}
