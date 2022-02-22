package storage

import "database/sql"

func NewDB(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return nil, err
	}

	return db, nil
}
