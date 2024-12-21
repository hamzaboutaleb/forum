package api

import (
	"net/http"
	"strconv"

	"forum/config"
	"forum/models"
	"forum/utils"
)

func ReactToPostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleReactPost(w, r)
	case http.MethodGet:
		handleReactGet(w, r)
	default:
		utils.WriteJSON(w, http.StatusMethodNotAllowed, "The HTTP method used in the request is invalid. Please ensure you're using the correct method.", nil)
	}
}

func handleReactPost(w http.ResponseWriter, r *http.Request) {
	sessionId := utils.GetSessionCookie(r)
	session := config.IsAuth(sessionId)
	if session == nil {
		utils.WriteJSON(w, http.StatusUnauthorized, "You don't have the necessary permissions to access this. Please log in or check your access rights.", nil)
		return
	}
	var like models.Like
	err := utils.ReadJSON(r, &like)
	like.UserID = session.UserId
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	postRepo := models.NewPostRepository()
	isExistPost, err := postRepo.IsPostExist(like.PostID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	if !isExistPost {
		utils.WriteJSON(w, http.StatusBadRequest, "The post you're trying to react to does not exist.", nil)
		return
	}
	if like.IsLike != -1 && like.IsLike != 1 {
		utils.WriteJSON(w, http.StatusBadRequest, "You can only like or dislike a post. Please choose one of these actions.", nil)
		return
	}
	likeRepo := models.NewLikeRepository()
	ok, err := likeRepo.IsUserReactToPost(like.UserID, like.PostID, like.IsLike)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	if ok {
		if err := likeRepo.DeleteLike(like.UserID, like.PostID); err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
	} else {
		if err := likeRepo.AddReaction(&like); err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
	}
	utils.WriteJSON(w, http.StatusOK, "like updated succefully", nil)
}

func handleReactGet(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.ParseInt(r.URL.Query().Get("postId"), 10, 64)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, "Bad Request", nil)
		return
	}
	likeRepo := models.NewLikeRepository()
	count, err := likeRepo.GetPostLikes(postId)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "", count)
}
