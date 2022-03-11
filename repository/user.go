package repository

import (
	"database/sql"
)

type User struct {
	Id           string         `json:"id"`
	Email        string         `json:"email"`
	Username     string         `json:"username"`
	Password     string         `json:"password"`
	Activated    bool           `json:"activated"`
	RefreshToken sql.NullString `json:"refresh_token"`
}

func (u *User) GetUserByEmail(db *sql.DB) error {
	return db.QueryRow(
		"SELECT * FROM users WHERE users.email = $1",
		u.Email).Scan(&u.Id, &u.Email, &u.Username, &u.Password, &u.RefreshToken, &u.Activated)
}

func (u *User) GetUserById(db *sql.DB) error {
	return db.QueryRow("SELECT * FROM users WHERE users.id = $1",
		u.Id).Scan(&u.Id, &u.Email, &u.Username, &u.Password, &u.RefreshToken, &u.Activated)
}

func (u *User) Activate(db *sql.DB) error {
	_, err := db.Exec("UPDATE users SET activated = TRUE WHERE id = $1", u.Id)

	return err
}

func (u *User) Create(db *sql.DB) error {
	_, err :=
		db.Exec(
			"INSERT INTO users (email, username, password) VALUES ($1, $2, $3)",
			u.Email, u.Username, u.Password,
		)

	return err
}
