package api

import (
	"encoding/json"
	"net/http"

	"forum/config"
	"forum/models"
	"forum/utils"
)

func ReactToPostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleReactGet(w, r)
	default:
		utils.WriteJSON(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
	}
}

func handleReactGet(w http.ResponseWriter, r *http.Request) {
	session := config.IsAuth(utils.GetSessionCookie(r))
	if session == nil {
		utils.WriteJSON(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	var like models.Like
	err := utils.ReadJSON(r, &like)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	like.UserID = session.UserId
	postRepo := models.NewPostRepository()
	isExistPost, _ := postRepo.IsPostExist(like.PostID)
	if isExistPost == 0 {
		utils.WriteJSON(w, http.StatusBadRequest, "Invalid Request", nil)
		return
	}
	if like.IsLike != -1 && like.IsLike != 1 {
		utils.WriteJSON(w, http.StatusBadRequest, "Invalid Request", nil)
		return
	}
	if err != nil {
	}

	likeRepo := models.NewLikeRepository()

	if err := likeRepo.AddReaction(&like); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Reaction added successfully"})
}
