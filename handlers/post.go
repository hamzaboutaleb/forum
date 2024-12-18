package handlers

import (
	"database/sql"
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
		config.TMPL.RenderError(w, "error.html", err.Error(), 500)
		return
	}
	postRepo := models.NewPostRepository()
	comRepo := models.NewCommentRepository()
	post, err := postRepo.GetPostById(postId)
	if err != nil {
		if err == sql.ErrNoRows {
			config.TMPL.RenderError(w, "error.html", "Not found", 404)
			return
		}
		config.TMPL.RenderError(w, "error.html", err.Error(), 500)
		return
	}
	comment, err := comRepo.GetPostComments(postId)
	if err != nil {
		config.TMPL.RenderError(w, "error.html", err.Error(), 500)
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
