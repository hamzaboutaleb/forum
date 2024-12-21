package handlers

import (
	"net/http"

	"forum/config"
	"forum/utils"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteJSON(w, http.StatusUnauthorized, "The HTTP method used in the request is invalid. Please ensure you're using the correct method.", nil)
		return
	}
	utils.RedirectIsAuth(w, r)
	config.TMPL.Render(w, "register.html", nil)
}
