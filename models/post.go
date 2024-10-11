package models

import (
	"database/sql"

	"forum/config"
)

type Post struct {
	ID        int64  `json:"id"`
	UserID    int    `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: config.DB}
}

func (r *PostRepository) Create(post *Post) error {
	query := `INSERT INTO posts (user_id, content, createdAt) VALUES (?, ?, ?)`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return config.NewInternalError(err)
	}
	result, err := stmt.Exec(post.UserID, post.Content, post.CreatedAt)
	if err != nil {
		return config.NewInternalError(err)
	}
	id, _ := result.LastInsertId()
	post.ID = id
	return nil
}
