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

type service struct {
	repo database.CommentRepo
}

func NewCommentService(repo database.CommentRepo) CommentService {
	return service{
		repo: repo,
	}
}

func (s service) SaveComment(comment domain.Comment) (domain.Comment, error) {
	return s.repo.SaveComment(comment)
}

func (s service) GetComment(id int64) (domain.Comment, error) {
	return s.repo.GetComment(id)
}

func (s service) UpdateComment(comment domain.Comment) (domain.Comment, error) {
	return s.repo.UpdateComment(comment)
}

func (s service) DeleteComment(id int64) error {
	return s.repo.DeleteComment(id)
}
