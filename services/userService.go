package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"forum/models"
)

var ErrUserOrEmailExist = errors.New("username or email already used")

func RegisterUser(user *models.User) error {
	userRepo := models.NewUserRepository()

	// check if the username or email alread yexist
	isUserExist, err := userRepo.UserExists(user.Username, user.Email)
	if err != nil {
		return err
	}
	if isUserExist {
		return ErrUserOrEmailExist
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return userRepo.CreateUser(user)
}
