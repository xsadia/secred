package storage

import "database/sql"

func GetDB(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return nil, err
	}

	return db, nil
}
