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
	GetPostsByUser(userID int64) ([]domain.Post, error)
}

type postService struct {
	repo    database.PostRepo
	service CommentService
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
	post, err := s.repo.GetPost(id)
	if err != nil {
		return domain.Post{}, err
	}
	post.Comment, err = s.service.GetCommentsByPostID(id)
	if err != nil {
		return domain.Post{}, err
	}
	return post, err
}

func (s postService) UpdatePost(post domain.Post) (domain.Post, error) {
	return s.repo.UpdatePost(post)
}

func (s postService) DeletePost(id int64) error {
	return s.repo.DeletePost(id)
}

func (s postService) GetPostsByUser(userID int64) ([]domain.Post, error) {
	return s.repo.GetPostsByUser(userID)
}
