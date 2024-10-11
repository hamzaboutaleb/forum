package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"forum/config"
	"forum/utils"
)

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
	session := utils.GeTCookie("session", r)
	page := NewPageStruct("forum", session, nil)
	config.TMPL.Render(w, "index.html", page)
}

func indexPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("post")
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
