package handlers

import (
	"net/http"

	"forum/config"
	"forum/services"
	"forum/utils"
)

type pageData struct {
	Error  string
	Method string
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	utils.RedirectIsAuth(w, r)

	page := pageData{
		Method: "POST",
	}
	if r.Method == http.MethodGet {
		page.Method = "GET"
		config.TMPL.Render(w, "login.html", page)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	if err := services.LoginUser(username, password); err != nil {
		// TODO make page
		page.Error = err.Error()
	}
	session, err := config.SESSION.CreateSession(username)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	cookies := http.Cookie{
		Name:    "session",
		Value:   session.ID,
		Expires: session.ExpiresAt,
		Path:    "/",
	}
	http.SetCookie(w, &cookies)
	config.TMPL.Render(w, "login.html", page)
}
