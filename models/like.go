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
	stmt, err := r.db.Prepare(`
        INSERT INTO post_reactions (userId, postId, isLike)
        VALUES (?, ?, ?)
        ON CONFLICT(userId, postId) DO UPDATE SET isLike = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(like.UserID, like.PostID, like.IsLike, like.IsLike)
	return err
}

func (r *LikeRepository) IsReactionExists(like *Like) (bool, int, error) {
	var exists bool
	var isLike int
	stmt, err := r.db.Prepare(`
        SELECT EXISTS(SELECT 1 FROM post_reactions WHERE userId = ? AND postId = ?),
               isLike
        FROM post_reactions WHERE userId = ? AND postId = ?`)
	if err != nil {
		return false, 0, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(like.UserID, like.PostID, like.UserID, like.PostID).Scan(&exists, &isLike)
	if err != nil && err != sql.ErrNoRows {
		return false, 0, err
	}

	return exists, isLike, nil
}

func (r *LikeRepository) CountLikes(postId int64) (int, error) {
	stmt, err := r.db.Prepare("SELECT SUM(isLike) FROM post_reactions WHERE postId = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var likes int
	err = stmt.QueryRow(postId).Scan(&likes)
	if err != nil {
		return 0, err
	}
	return likes, nil
}

func (r *LikeRepository) GetPostLikes(postId int64) (*PostLike, error) {
	stmt, err := r.db.Prepare("select postId, SUM(CASE WHEN isLike = -1 THEN 1 ELSE 0 END) as dislike, SUM(CASE WHEN IsLike = 1 THEN 1 ELSE 0 END) as likes from post_reactions WHERE postId = ? GROUP BY postId")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var postLike PostLike
	err = stmt.QueryRow(postId).Scan(&postLike.PostId, &postLike.DislikesCount, &postLike.LikesCount)
	if err != nil {
		return nil, err
	}
	return &postLike, nil
}
