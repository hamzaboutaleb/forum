package api

import (
	"fmt"
	"log"
	"net/http"

	"forum/config"
	"forum/models"
	"forum/utils"
)

func AddComment(w http.ResponseWriter, r *http.Request) {
	sessionId := utils.GetSessionCookie(r)
	session := config.IsAuth(sessionId)
	fmt.Println("id", sessionId)
	if session == nil {
		utils.WriteJSON(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	if r.Method != http.MethodPost {
		utils.WriteJSON(w, http.StatusUnauthorized, "Invalid request method", nil)
		return
	}

	var comment models.Comment
	comment.UserID = session.UserId
	err := utils.ReadJSON(r, &comment)
	fmt.Println(comment.Comment)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if comment.PostID == 0 || comment.Comment == "" {
		http.Error(w, "comment is required", http.StatusBadRequest)
		return
	}

	commRepo := models.NewCommentRepository()
	commRepo.Create(&comment)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	utils.WriteJSON(w, 200, "comment added succefully", comment)
}

// userId, postId, isLike
func HandleLikeComment(w http.ResponseWriter, r *http.Request) {
	sessionId := utils.GetSessionCookie(r)
	session := config.IsAuth(sessionId)
	if session == nil {
		utils.WriteJSON(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	var like models.CommentLike
	like.UserID = session.UserId
	err := utils.ReadJSON(r, &like)
	fmt.Println(like.CommentId, like.UserID)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	if like.IsLike != 1 && like.IsLike != -1 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	comntRepo := models.NewCommentRepository()
	err = comntRepo.ReactComment(like)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, 200, "comment added succefully", like)
}
