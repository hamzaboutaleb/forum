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
	switch r.Method {
	case http.MethodGet:
		getLogin(w)
	case http.MethodPost:
		postLogin(w, r)
	default:
		config.TMPL.RenderError(w, "error.html", "The HTTP method used in the request is invalid. Please ensure you're using the correct method.", http.StatusMethodNotAllowed)
	}
}

func getLogin(w http.ResponseWriter) {
	page := pageData{
		Method: "GET",
	}
	config.TMPL.Render(w, "login.html", page)
}

func postLogin(w http.ResponseWriter, r *http.Request) {
	page := pageData{
		Method: "POST",
	}
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	user, err := services.LoginUser(username, password)
	if err != nil {
		page.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		config.TMPL.Render(w, "login.html", page)
		return
	}

	session, err := config.SESSION.CreateSession(user.Username, user.ID)
	if err != nil {
		config.TMPL.RenderError(w, "error.html", err.Error(), http.StatusInternalServerError)
		return
	}
	cookies := http.Cookie{
		Name:    "session",
		Value:   session.ID,
		Expires: session.ExpiresAt,
		Path:    "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookies)
	config.TMPL.Render(w, "login.html", page)
}
