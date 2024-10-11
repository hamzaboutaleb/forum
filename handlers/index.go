package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"forum/config"
	"forum/models"
	"forum/utils"
)

type IndexStruct struct {
	Posts     []models.Post
	PostCount int
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Page Not Found", http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		indexGet(w, r)
	case http.MethodPost:
		indexPost(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func indexGet(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	currPage, err := strconv.Atoi(pageStr)
	if err != nil || currPage < 1 {
		currPage = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = config.LIMIT_PER_PAGE
	}
	session := utils.GeTCookie("session", r)
	postRep := models.NewPostRepository()
	post, err := postRep.GetPostPerPage(currPage, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	count, err := postRep.Count()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	page := NewPageStruct("forum", session, nil)
	page.Data = IndexStruct{
		Posts:     *post,
		PostCount: count,
	}
	config.TMPL.Render(w, "index.html", page)
}

func indexPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	session := utils.GeTCookie("session", r)
	page := NewPageStruct("forum", session, nil)
	title := r.FormValue("title")
	content := r.FormValue("content")
	tags := r.FormValue("tags")
	response := Response{}
	if strings.TrimSpace(title) == "" || strings.TrimSpace(content) == "" || strings.TrimSpace(tags) == "" {
		response.Error = true
		response.Message = "All fields must be completed."
	}
	page.Data = response
	// TODO chekc tags and insert them
	config.TMPL.Render(w, "index.html", page)
}
