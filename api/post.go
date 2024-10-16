package api

import (
	"fmt"
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
		utils.WriteJSON(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	sessionId := utils.GeTCookie("session", r)
	session, err := config.SESSION.GetSession(sessionId)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	var post models.Post
	err = utils.ReadJSON(r, &post)
	if err != nil {
		fmt.Println("here")
		utils.WriteJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	post.UserID = session.UserId
	err = services.CreateNewPost(&post)
	if err != nil {
		if err.(*config.CustomError).IsInternal() {
			utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		utils.WriteJSON(w, http.StatusBadRequest, err.Error(), nil)
	}
	utils.WriteJSON(w, http.StatusCreated, "Post created successfully", post)
}
