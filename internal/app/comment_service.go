package app

import (
	"trainee/internal/domain"
	"trainee/internal/infra/database"
)

type CommentService interface {
	SaveComment(comment domain.Comment) (domain.Comment, error)
	GetComment(id int64) (domain.Comment, error)
	UpdateComment(comment domain.Comment) (domain.Comment, error)
	DeleteComment(id int64) error
}

type commentService struct {
	repo database.CommentRepo
}

func NewCommentService(repo database.CommentRepo) CommentService {
	return commentService{
		repo: repo,
	}
}

func (s commentService) SaveComment(comment domain.Comment) (domain.Comment, error) {
	return s.repo.SaveComment(comment)
}

func (s commentService) GetComment(id int64) (domain.Comment, error) {
	return s.repo.GetComment(id)
}

func (s commentService) UpdateComment(comment domain.Comment) (domain.Comment, error) {
	return s.repo.UpdateComment(comment)
}

func (s commentService) DeleteComment(id int64) error {
	return s.repo.DeleteComment(id)
}
