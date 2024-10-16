package models

import (
	"database/sql"
	"fmt"
	"time"

	"forum/config"
)

type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"postId"`
	UserID    int64     `json:"userId"`
	Username  string    `json:"username"`
	Likes     int       `json:"likes"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
}

type CommentRepository struct {
	db *sql.DB
}

func NewCommnetRepository() *CommentRepository {
	return &CommentRepository{db: config.DB}
}

func (r *CommentRepository) Create(comment *Comment) error {
	query := `INSERT INTO comments (postId, userId, comment) VALUES (?, ?, ?)`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(comment.PostID, comment.UserID, comment.Comment)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	comment.ID = id
	return nil
}

func (r *CommentRepository) GetPostComments(postID int64) ([]Comment, error) {
	query := `SELECT c.id ,c.postId ,c.userId, c.comment ,c.createdAt ,u.username ,COALESCE(SUM(l.isLike), 0) AS likeCount FROM comments c 
	LEFT JOIN comment_reactions l ON c.id = l.commentId 
	LEFT JOIN users u ON c.userId = u.id WHERE c.postId = ? GROUP BY c.id HAVING count(c.id) > 0`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var comments []Comment
	rows, err := stmt.Query(postID)
	if err != nil {
		return comments, nil
	}
	defer rows.Close()

	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Comment, &comment.CreatedAt,
			&comment.Username, &comment.Likes)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}
