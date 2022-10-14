package app

import (
	"trainee/internal/domain"
	"trainee/internal/infra/database"
)

type CommentService interface {
	SaveComment(comment domain.Comment) (domain.Comment, error)
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
