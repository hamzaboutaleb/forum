package api

import (
	"net/http"

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
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, "Invalid input", nil)
		return
	}
	if !utils.IsValidEmail(user.Email) {
		utils.WriteJSON(w, http.StatusBadRequest, "Invalid Email", nil)
		return
	}
	err = services.RegisterUser(&user)
	if err != nil {
		if err == services.ErrUserOrEmailExist {
			utils.WriteJSON(w, http.StatusConflict, "Username or Email already used", nil)
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, "User registered successfully", user)
}
