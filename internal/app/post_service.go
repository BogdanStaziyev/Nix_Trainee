package app

import (
	"trainee/internal/domain"
	"trainee/internal/infra/database"
)

//go:generate mockery --dir . --name PostService --output ./mocks
type PostService interface {
	SavePost(post domain.Post) (domain.Post, error)
	GetPost(id int64) (domain.Post, error)
	UpdatePost(post domain.Post) (domain.Post, error)
	DeletePost(id int64) error
}

type postService struct {
	repo database.PostRepo
}

func NewPostService(repo database.PostRepo) PostService {
	return postService{
		repo: repo,
	}
}

func (s postService) SavePost(post domain.Post) (domain.Post, error) {
	return s.repo.SavePost(post)
}

func (s postService) GetPost(id int64) (domain.Post, error) {
	return s.repo.GetPost(id)
}

func (s postService) UpdatePost(post domain.Post) (domain.Post, error) {
	return s.repo.UpdatePost(post)
}

func (s postService) DeletePost(id int64) error {
	return s.repo.DeletePost(id)
}
