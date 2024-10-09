package handlers

import (
	"net/http"

	"forum/config"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	config.TMPL.Render(w, "register.html", nil)
}
