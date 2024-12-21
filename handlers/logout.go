package handlers

import (
	"net/http"
	"time"

	"forum/config"
	"forum/utils"
)

func deleteCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func LogutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteJSON(w, http.StatusUnauthorized, "The HTTP method used in the request is invalid. Please ensure you're using the correct method.", nil)
		return
	}
	sessionId := utils.GetSessionCookie(r)
	config.SESSION.DeleteSession(sessionId)
	deleteCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
