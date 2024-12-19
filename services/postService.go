package services

import (
	"errors"
	"strings"
	"time"

	"forum/config"
	"forum/models"
)

func CheckTags(strs []string) bool {
	for _, str := range strs {
		if len(str) > 20 {
			return false
		}
	}
	return true
}

func IsBetween(str string, x int, y int) bool {
	if len(str) >= x && len(str) <= y {
		return true
	}
	return false
}

func CreateNewPost(post *models.Post) error {
	postRepo := models.NewPostRepository()
	TagsRepo := models.NewTagRepository()
	if !IsBetween(post.Title, 0, 200) {
		return config.NewError(errors.New("title has exceeded the limits"))
	}
	if !IsBetween(post.Content, 0, 3000) {
		return config.NewError(errors.New("content has exceeded the limits"))
	}
	if !CheckTags(post.Tags) {
		return config.NewError(errors.New("tags have exceeded the limits"))
	}
	// check if input empty
	if strings.TrimSpace(post.Content) == "" || post.IsTagsEmpty() {
		return config.NewError(errFieldsEmpty)
	}
	post.CreatedAt = time.Now()
	err := postRepo.Create(post)
	if err != nil {
		return err
	}
	TagsRepo.LinkTagsToPost(post.ID, post.Tags)
	return err
}
