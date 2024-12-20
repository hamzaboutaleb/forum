package models

import (
	"database/sql"

	"forum/config"
)

type Like struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"userId"`
	PostID int64 `json:"postid"`
	IsLike int   `json:"isLike"`
}

type PostLike struct {
	PostId        int64 `json:"postId"`
	LikesCount    int64 `json:"likesCount"`
	DislikesCount int64 `json:"dislikesCount"`
}

type LikeRepository struct {
	db *sql.DB
}

func NewLikeRepository() *LikeRepository {
	return &LikeRepository{db: config.DB}
}

func (r *LikeRepository) AddReaction(like *Like) error {
	query := `
        INSERT INTO post_reactions (userId, postId, isLike)
        VALUES (?, ?, ?)
        ON CONFLICT(userId, postId) DO UPDATE SET isLike = ?`
	_, err := r.db.Exec(query, like.UserID, like.PostID, like.IsLike, like.IsLike)
	return err
}

func (r *LikeRepository) IsReactionExists(like *Like) (bool, int, error) {
	var exists bool
	var isLike int
	query := `
        SELECT EXISTS(SELECT 1 FROM post_reactions WHERE userId = ? AND postId = ?),
               isLike
        FROM post_reactions WHERE userId = ? AND postId = ?`
	err := r.db.QueryRow(query, like.UserID, like.PostID, like.UserID, like.PostID).Scan(&exists, &isLike)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, 0, nil
		}
		return false, 0, err
	}

	return exists, isLike, nil
}

func (r *LikeRepository) CountLikes(postId int64) (int, error) {
	query := "SELECT SUM(isLike) FROM post_reactions WHERE postId = ?"
	var likes int
	err := r.db.QueryRow(query, postId).Scan(&likes)
	if err != nil {
		return 0, err
	}
	return likes, nil
}

func (r *LikeRepository) GetPostLikes(postId int64) (*PostLike, error) {
	query := `select 
	postId, 
	SUM(CASE WHEN isLike = -1 THEN 1 ELSE 0 END) as dislike, 
	SUM(CASE WHEN IsLike = 1 THEN 1 ELSE 0 END) as likes 
	from post_reactions WHERE postId = ? 
	GROUP BY postId`

	var postLike PostLike
	err := r.db.QueryRow(query, postId).Scan(&postLike.PostId, &postLike.DislikesCount, &postLike.LikesCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return &postLike, nil
		}
		return nil, err
	}
	return &postLike, nil
}

func (r *LikeRepository) IsUserReactToPost(userId int64, postId int64, isLike int) (bool, error) {
	query := `SELECT COUNT(id) FROM post_reactions WHERE userId = ? AND postId = ? AND isLike = ?`
	var numRows int

	err := r.db.QueryRow(query, userId, postId, isLike).Scan(&numRows)
	if err != nil {
		return false, err
	}
	return numRows > 0, nil
}

func (r *LikeRepository) DeleteLike(userId int64, postId int64) error {
	query := "DELETE FROM post_reactions WHERE userId=? AND postId=?"
	_, err := r.db.Exec(query, userId, postId)
	if err != nil {
		return err
	}
	return nil
}
