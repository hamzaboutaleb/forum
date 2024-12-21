package api

import (
	"net/http"
	"strings"

	"forum/config"
	"forum/models"
	"forum/utils"
)

func AddComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSON(w, http.StatusUnauthorized, "The HTTP method used in the request is invalid. Please ensure you're using the correct method.", nil)
		return
	}
	sessionId := utils.GetSessionCookie(r)
	session := config.IsAuth(sessionId)
	if session == nil {
		utils.WriteJSON(w, http.StatusUnauthorized, "You don't have the necessary permissions to access this. Please log in or check your access rights.", nil)
		return
	}

	var comment models.Comment
	err := utils.ReadJSON(r, &comment)
	comment.UserID = session.UserId
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if strings.TrimSpace(comment.Comment) == "" || !utils.IsBetween(comment.Comment, 2, 1000) {
		utils.WriteJSON(w, http.StatusBadRequest, "The comment must be between 2 and 1000 characters", nil)
		return
	}
	postRepo := models.NewPostRepository()
	commentRepo := models.NewCommentRepository()
	isExist, err := postRepo.IsPostExist(comment.PostID)
	if err != nil {
		config.TMPL.RenderError(w, "error.html", "Internal server error", http.StatusInternalServerError)
		return
	}
	if !isExist {
		utils.WriteJSON(w, http.StatusBadRequest, "It looks like the post you're trying to comment on doesn't exist anymore.", nil)
		return
	}
	err = commentRepo.Create(&comment)
	if err != nil {
		config.TMPL.RenderError(w, "error.html", "Internal server error", http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, 200, "Your comment has been added successfully! Thanks for sharing your thoughts!", comment)
}


func HandleLikeComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSON(w, http.StatusUnauthorized, "The HTTP method used in the request is invalid. Please ensure you're using the correct method.", nil)
		return
	}
	sessionId := utils.GetSessionCookie(r)
	session := config.IsAuth(sessionId)
	if session == nil {
		utils.WriteJSON(w, http.StatusUnauthorized, "You don't have the necessary permissions to access this. Please log in or check your access rights.", nil)
		return
	}
	var like models.CommentLike
	err := utils.ReadJSON(r, &like)
	like.UserID = session.UserId
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	if like.IsLike != 1 && like.IsLike != -1 {
		utils.WriteJSON(w, http.StatusBadRequest, "You can only like or dislike a comment. Please choose one of these actions.", nil)
		return
	}
	comntRepo := models.NewCommentRepository()
	isExist, err := comntRepo.IsCommentExist(like.CommentId)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	if !isExist {
		utils.WriteJSON(w, http.StatusBadRequest, "You can't like or dislike a comment that doesn't exist. It might have been removed.", nil)
		return
	}
	isReactionExist, err := comntRepo.IsReactionExist(like.UserID, like.CommentId, like.IsLike)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	if isReactionExist {
		err = comntRepo.DeleteReaction(like.UserID, like.CommentId)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
	} else {
		err = comntRepo.ReactComment(like)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
	}
	commentReaction, err := comntRepo.GetCommentReaction(like.CommentId)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "like updated succefully.", commentReaction)
}
