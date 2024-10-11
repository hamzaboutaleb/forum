package api

import (
	"net/http"

	"forum/utils"
)

func NewPostApi(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlePost(w, r)
	default:
		utils.WriteJSON(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
}
