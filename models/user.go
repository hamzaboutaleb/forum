package models

import (
	"database/sql"
	"errors"
	"fmt"

	"forum/config"

	"github.com/gofrs/uuid/v5"
)

var (
	ErrDB           = errors.New("database error")
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	ID       string `json:"id"`
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
	query := "INSERT INTO users (id, email, username, password) VALUES (?,?,?,?)"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	userID, err := uuid.NewV7()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userID.String(), user.Email, user.Username, user.Password)
	if err != nil {
		return err
	}
	user.ID = userID.String()
	return nil
}

func (r *UserRepository) GetUserByID(id string) (*User, error) {
	query := "SELECT id, email, username, password FROM users WHERE id = ?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, ErrDB
	}
	row := stmt.QueryRow(id)
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
		return nil, ErrDB
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
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, ErrDB
	}
	row := stmt.QueryRow(username)

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
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return false, err
	}
	err = stmt.QueryRow(username, email).Scan(&count)
	fmt.Println(username, email, count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
