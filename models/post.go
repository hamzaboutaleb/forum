package models

import (
	"database/sql"
	"strings"
	"time"

	"forum/config"
)

const (
	ALL = iota
	MY_POST
	LIKED_POST
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
	Dislikes  int       `json:"dislikes"`
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

	result, err := r.db.Exec(query, post.UserID, post.Title, post.Content, post.CreatedAt)
	if err != nil {
		return err
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
		SUM(CASE WHEN isLike = 1 THEN 1 ELSE 0 END) as likes,
		SUM(CASE WHEN isLike = -1 THEN 1 ELSE 0 END) as dislike
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

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	var posts []*Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.Username, &post.Likes, &post.Dislikes); err != nil {
			if err == sql.ErrNoRows {
				return posts, nil
			}
			return nil, err
		}
		posts = append(posts, &post)
	}
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
    SUM(CASE WHEN pr.isLike = 1 THEN 1 ELSE 0 END) as likes,
		SUM(CASE WHEN pr.isLike = -1 THEN 1 ELSE 0 END) as dislike
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

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.Username, &post.UserID, &post.Likes, &post.Dislikes); err != nil {
			if err == sql.ErrNoRows {
				return posts, nil
			}
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *PostRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM posts`).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}

func (r *PostRepository) IsPostExist(id int64) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM posts WHERE id = ?`
	err := r.db.QueryRow(query, id).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return count > 0, nil
}

func (r *PostRepository) GetPostById(id int64) (*Post, error) {
	query := `SELECT * FROM posts WHERE id = ?`

	var post Post
	row := r.db.QueryRow(query, id)
	err := row.Scan(&post.ID, &post.Title, &post.UserID, &post.Content, &post.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) GetPostsByTag(tag string) ([]Post, error) {
	query := `SELECT p.id, p.title, p.content, p.createdAt, u.username, COALESCE(SUM(pl.isLike), 0) AS likeCount 
	FROM posts p
	LEFT JOIN users u ON p.userId = u.id
	LEFT JOIN post_reactions pr ON p.id = pr.postId
	GROUP BY p.id 
	HAVING p.id IN (SELECT postId from post_tags WHERE tagId = (SELECT id FROM tags WHERE name=?))`
	posts := []Post{}

	rows, err := r.db.Query(query, tag)
	if err != nil {
		return posts, err
	}
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.Username, &post.Likes)
		if err != nil {
			if err == sql.ErrNoRows {
				return posts, nil
			}
			return posts, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}


func (r *PostRepository) CompleteQuery(query, tag string, queryType int, userId int64, page int, limit int) (*sql.Rows, error) {
	querys := []string{}
	prepare := []any{}
	if tag != "" {
		querys = append(querys, "p.id IN (SELECT postId from post_tags WHERE tagId = (SELECT id FROM tags WHERE name=?))")
		prepare = append(prepare, tag)
	}
	switch queryType {
	case MY_POST:
		{
			querys = append(querys, "p.id IN (SELECT id FROM posts WHERE userId = ?)")
			prepare = append(prepare, userId)
		}
	case LIKED_POST:
		{
			querys = append(querys, "p.id IN (SELECT postId FROM post_reactions WHERE userId = ?)")
			prepare = append(prepare, userId)
		}
	}
	if len(querys) > 0 {
		querys[0] = " HAVING " + querys[0]
	}
	queryStr := query + strings.Join(querys, " AND ")
	queryStr += " ORDER BY p.createdAt DESC"

	rows, err1 := r.db.Query(queryStr, prepare...)
	return rows, err1
}

func (r *PostRepository) GetPostsBy(tag string, filterType int, userId int64, page, limit int) ([]*Post, error) {
	queryPostIds := `SELECT p.id, p.title, p.content, p.createdAt, u.username,
	SUM(CASE WHEN pr.isLike = 1 THEN 1 ELSE 0 END) as likes,
	SUM(CASE WHEN pr.isLike = -1 THEN 1 ELSE 0 END) as dislike
	FROM posts p
	LEFT JOIN users u ON p.userId = u.id
	LEFT JOIN post_reactions pr ON p.id = pr.postId
	GROUP BY p.id
	`
	posts := []*Post{}

	rows, err := r.CompleteQuery(queryPostIds, tag, filterType, userId, page, limit)
	if err != nil {
		return posts, err
	}
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.Username, &post.Likes, &post.Dislikes)
		if err != nil {
			if err == sql.ErrNoRows {
				return posts, nil
			}
			return posts, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}
