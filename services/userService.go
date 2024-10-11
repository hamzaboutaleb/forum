package services

import (
	"errors"
	"strings"

	"forum/config"
	"forum/models"
	"forum/utils"
)

var (
	ErrInvalidUserPass  = errors.New("invalid username or password")
	ErrUserOrEmailExist = errors.New("username or email already used")
	ErrFieldsEmpty      = errors.New("all fields must be completed")
)

func RegisterUser(user *models.User) error {
	userRepo := models.NewUserRepository()

	// check if the username or email alread yexist
	isUserExist, err := userRepo.UserExists(user.Username, user.Email)
	if err != nil {
		return err
	}
	if isUserExist {
		return config.NewError(ErrUserOrEmailExist)
	}
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return config.NewInternalError(err)
	}
	user.Password = hashedPassword
	return userRepo.CreateUser(user)
}

func LoginUser(username, password string) error {
	if len(strings.TrimSpace(username)) == 0 || len(strings.TrimSpace(password)) == 0 {
		return config.NewError(ErrFieldsEmpty)
	}
	userRepo := models.NewUserRepository()
	user, err := userRepo.GetUserByUsername(username)
	if err != nil {
		return err
	}
	if utils.CheckPassword(user.Password, password) != nil {
		return config.NewError(ErrInvalidUserPass)
	}
	return nil
}
