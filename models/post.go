package models

import (
	"database/sql"
	"strings"

	"forum/config"
)

type Post struct {
	ID        int64    `json:"id"`
	UserID    string   `json:"userId"`
	Content   string   `json:"content"`
	CreatedAt string   `json:"createdAt"`
	Tags      []string `json:"tags"`
}

type PostRepository struct {
	db *sql.DB
}

func (p *Post) IsTagsEmpty() bool {
	for _, tag := range p.Tags {
		if strings.TrimSpace(tag) == "" {
			return true
		}
	}
	return false
}

func NewPostRepository() *PostRepository {
	return &PostRepository{db: config.DB}
}

func (r *PostRepository) Create(post *Post) error {
	query := `INSERT INTO posts (userId, content, createdAt) VALUES (?, ?, ?)`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return config.NewInternalError(err)
	}
	defer stmt.Close()
	result, err := stmt.Exec(post.UserID, post.Content, post.CreatedAt)
	if err != nil {
		return config.NewInternalError(err)
	}
	id, _ := result.LastInsertId()
	post.ID = id
	return nil
}

func (r *PostRepository) FindAll() ([]Post, error) {
	query := "SELECT * FROM posts"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, config.NewInternalError(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil, config.NewInternalError(err)
	}
	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.CreatedAt); err != nil {
			if err == sql.ErrNoRows {
				return nil, config.NewError(err)
			}
			return nil, config.NewInternalError(err)
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, config.NewInternalError(err)
	}
	return posts, nil
}
