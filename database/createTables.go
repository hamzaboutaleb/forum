package database

import (
	"forum/config"
)

func CreateUserTable() error {
	db := config.DB
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
