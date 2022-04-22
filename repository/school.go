package repository

import "database/sql"

type School struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (s *School) CreateSchool(db *sql.DB) error {
	return db.QueryRow("INSERT INTO schools (name) VALUES ($1) RETURNING id", s.Name).Scan(&s.Id)
}
