package database

import (
	"forum/config"
)

func execQuery(query string) error {
	db := config.DB
	_, err := db.Exec(query)
	if err != nil {
		return config.NewInternalError(err)
	}
	return nil
}

var Tables = []func() error{
	createPostTable,
	createUserTable,
	createSessionTable,
}

func createUserTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);
	`
	return execQuery(query)
}

func createSessionTable() error {
	query := `CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		username TEXT,
		userId TEXT,
		expires_at DATETIME
	)`

	return execQuery(query)
}

func createPostTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS posts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
        userId INTEGER NOT NULL,
        content TEXT NOT NULL,
        createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	return execQuery(query)
}
