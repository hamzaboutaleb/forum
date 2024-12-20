package models

import (
	"database/sql"
	"errors"

	"forum/config"
)

var (
	ErrDB           = errors.New("database error")
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db: config.DB}
}

func (r *UserRepository) CreateUser(user *User) error {
	query := "INSERT INTO users (email, username, password) VALUES (?,?,?)"
	result, err := r.db.Exec(query, user.Email, user.Username, user.Password)
	if err != nil {
		return err
	}
	user.ID, _ = result.LastInsertId()
	return nil
}

func (r *UserRepository) GetUserByID(id string) (*User, error) {
	query := "SELECT id, email, username, password FROM users WHERE id = ?"
	row := r.db.QueryRow(query, id)
	var user User
	if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	query := "SELECT id, email, username, password FROM users WHERE email = ?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(email)

	var user User
	if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*User, error) {
	query := "SELECT id, email, username, password FROM users WHERE username = ?"
	row := r.db.QueryRow(query, username)
	var user User
	if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UserExists(username, email string) (bool, error) {
	var count int
	query := `
    SELECT COUNT(*) FROM users 
    WHERE username = ? OR email = ?
    `
	err := r.db.QueryRow(query, username, email).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, ErrUserNotFound
		}
		return false, err
	}
	return count > 0, nil
}
