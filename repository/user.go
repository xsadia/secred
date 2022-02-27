package repository

import (
	"database/sql"
	"errors"
)

type User struct {
	Id           string `json:"id"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token"`
}

func (u *User) GetUserByEmail(db *sql.DB) error {
	return db.QueryRow(
		"SELECT id, email, username FROM users WHERE users.email = $1", u.Email,
	).Scan(&u.Id, &u.Email, &u.Username)
}

func (u *User) Create(db *sql.DB) error {
	err := u.GetUserByEmail(db)

	if err != sql.ErrNoRows {
		return errors.New("e-mail already in use")
	}

	_, err =
		db.Exec(
			"INSERT INTO users (email, username, password) VALUES ($1, $2, $3)",
			u.Email, u.Username, u.Password,
		)

	return err
}
