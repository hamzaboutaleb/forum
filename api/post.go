package api

import (
	"net/http"

	"forum/config"
	"forum/models"
	"forum/services"
	"forum/utils"
)

func PostApi(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlePost(w, r)
	default:
		utils.WriteJSON(w, http.StatusMethodNotAllowed, "The HTTP method used in the request is invalid. Please ensure you're using the correct method.", nil)
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	sessionId := utils.GetSessionCookie(r)
	session, _ := config.SESSION.GetSession(sessionId)
	if session == nil {
		utils.WriteJSON(w, http.StatusBadRequest, "You don't have the necessary permissions to access this. Please log in or check your access rights.", nil)
		return
	}
	var post models.Post
	err := utils.ReadJSON(r, &post)
	post.UserID = session.UserId
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	if utils.IsEmpty(post.Content) || utils.IsEmpty(post.Title) {
		utils.WriteJSON(w, http.StatusBadRequest, "The title or content cannot be empty. Please provide both and try again.", nil)
		return
	}
	err = services.CreateNewPost(&post)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, "The post has been created successfully.", post)
}
