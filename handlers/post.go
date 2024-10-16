package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"forum/config"
	"forum/models"
	"forum/utils"
)

type PostData struct {
	Post     models.Post
	Comments []models.Comment
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Println(postId)
	postRepo := models.NewPostRepository()
	comRepo := models.NewCommnetRepository()
	post, err := postRepo.GetPostById(postId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	comment, err := comRepo.GetPostComments(postId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	postData := PostData{
		Comments: comment,
		Post:     *post,
	}
	session := utils.GeTCookie("session", r)
	page := NewPageStruct(post.Title, session, postData)
	config.TMPL.Render(w, "post.html", page)
}