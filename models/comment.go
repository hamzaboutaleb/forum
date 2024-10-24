package models

import (
	"database/sql"
	"time"

	"forum/config"
)

type CommentLike struct {
	ID        int64 `json:"id"`
	UserID    int64 `json:"userId"`
	CommentId int64 `json:"commentId"`
	IsLike    int   `json:"isLike"`
}

type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"postId"`
	UserID    int64     `json:"userId"`
	Username  string    `json:"username"`
	Likes     int       `json:"likes"`
	DisLikes  int       `json:"disLikes"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
}

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository() *CommentRepository {
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
	query := `SELECT c.id ,c.postId ,c.userId, c.comment ,c.createdAt ,u.username , 
	(SELECT count(*) from comment_reactions WHERE isLike=1 AND commentId=c.id ) likes,
	(SELECT count(*) from comment_reactions WHERE isLike=-1 AND commentId=c.id ) dislike
	FROM comments c 
	LEFT JOIN comment_reactions l ON c.id = l.commentId 
	LEFT JOIN users u ON c.userId = u.id WHERE c.postId = ? GROUP BY c.id HAVING count(c.id) > 0 ORDER BY c.createdAt desc`
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
			&comment.Username, &comment.Likes, &comment.DisLikes)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *CommentRepository) ReactComment(like CommentLike) error {
	stmt, err := r.db.Prepare(`
        INSERT INTO comment_reactions (userId, commentId, isLike)
        VALUES (?, ?, ?)
        ON CONFLICT(userId, commentId) DO UPDATE SET isLike = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(like.UserID, like.CommentId, like.IsLike, like.IsLike)
	return err
}
