package services

import (
	"strings"

	"forum/config"
	"forum/models"
)

func CreateNewPost(post *models.Post) error {
	postRepo := models.NewPostRepository()
	// check if input empty
	if strings.TrimSpace(post.Content) == "" || post.IsTagsEmpty() {
		return config.NewError(errFieldsEmpty)
	}
	return postRepo.Create(post)
}
