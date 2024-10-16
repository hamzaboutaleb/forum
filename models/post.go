package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"forum/config"
)

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	UserID    int64     `json:"userId"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	Tags      []string  `json:"tags"`
	Username  string    `json:"Username"`
	Likes     int       `json:"likes"`
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
	query := `INSERT INTO posts (userId, title, content, createdAt) VALUES (?,?, ?, ?)`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return config.NewInternalError(err)
	}
	defer stmt.Close()
	result, err := stmt.Exec(post.UserID, post.Title, post.Content, post.CreatedAt)
	if err != nil {
		return config.NewInternalError(err)
	}
	id, _ := result.LastInsertId()
	post.ID = id
	return nil
}

func (r *PostRepository) GetPostPerPage(page int, limit int) ([]*Post, error) {
	offset := (page - 1) * limit
	query := `SELECT 
    p.id,
    p.title,
    p.content,
    p.createdAt,
    u.username,
	COALESCE(SUM(pl.isLike), 0) AS likeCount
	FROM 
    	posts p
	LEFT JOIN 
   		users u ON p.userId = u.id 
	LEFT JOIN 
    	post_reactions pl ON p.id = pl.postId
	GROUP BY 
    	p.id, u.id
	ORDER BY 
    	p.createdAt DESC
	LIMIT ? OFFSET ?`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, config.NewInternalError(err)
	}
	rows, err := stmt.Query(limit, offset)
	if err != nil {
		return nil, config.NewError(err)
	}
	var posts []*Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.Username, &post.Likes); err != nil {
			if err == sql.ErrNoRows {
				return nil, config.NewError(err)
			}
			return nil, config.NewInternalError(err)
		}
		posts = append(posts, &post)
	}
	defer stmt.Close()
	return posts, nil
}

func (r *PostRepository) FindAll() ([]Post, error) {
	query := `SELECT 
    p.id,
    p.title,
    p.content,
    p.createdAt,
    u.username,
	u.id
    SUM(pr.id) AS likeCount
	FROM 
    	posts p
	LEFT JOIN 
   		users u ON p.userId = u.id 
	LEFT JOIN 
    	post_reactions pr ON p.id = pr.postId
	GROUP BY 
    	p.id, u.id
	ORDER BY 
    	p.createdAt DESC`
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
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.Username, &post.UserID, &post.Likes); err != nil {
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
	fmt.Println(posts[0])
	return posts, nil
}

func (r *PostRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM posts`).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *PostRepository) IsPostExist(id int64) (int, error) {
	var count int
	stmt, err := r.db.Prepare(`SELECT COUNT(*) FROM posts WHERE id = ?`)
	if err != nil {
		return 0, err
	}
	err = stmt.QueryRow(id).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *PostRepository) GetPostById(id int64) (*Post, error) {
	query := `SELECT * FROM posts WHERE id = ?`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, config.NewInternalError(err)
	}
	var post Post
	row := stmt.QueryRow(id)
	err = row.Scan(&post.ID, &post.Title, &post.UserID, &post.Content, &post.CreatedAt)
	if err != nil {
		return nil, config.NewInternalError(err)
	}
	return &post, nil
}
