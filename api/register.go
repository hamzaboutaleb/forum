package api

import (
	"net/http"
	"strings"

	"forum/models"
	"forum/services"
	"forum/utils"
)

func RegisterApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}
	var user models.User
	err := utils.ReadJSON(r, &user)
	if !services.IsBetween(user.Username, 3, 50) {
		utils.WriteJSON(w, http.StatusBadRequest, "Invalid username", nil)
		return
	}
	if !services.IsBetween(user.Password, 3, 50) {
		utils.WriteJSON(w, http.StatusBadRequest, "Invalid password", nil)
		return
	}
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, "Invalid input", nil)
		return
	}
	if !utils.IsValidEmail(user.Email) {
		utils.WriteJSON(w, http.StatusBadRequest, "Invalid Email", nil)
		return
	}
	if len(strings.TrimSpace(user.Password)) == 0 {
		utils.WriteJSON(w, http.StatusBadRequest, "Invalid password", nil)
		return
	}
	err = services.RegisterUser(&user)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, "User registered successfully", user)
}
