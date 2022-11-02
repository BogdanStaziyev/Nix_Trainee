package app

import (
	"trainee/internal/domain"
	"trainee/internal/infra/database"
)

//go:generate mockery --dir . --name CommentService --output ./mocks
type CommentService interface {
	SaveComment(comment domain.Comment) (domain.Comment, error)
	GetComment(id int64) (domain.Comment, error)
	UpdateComment(comment string, id int64) (domain.Comment, error)
	DeleteComment(id int64) error
	GetCommentsByPostID(postID int64) ([]domain.Comment, error)
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

func (s commentService) UpdateComment(comment string, id int64) (domain.Comment, error) {
	return s.repo.UpdateComment(comment, id)
}

func (s commentService) DeleteComment(id int64) error {
	return s.repo.DeleteComment(id)
}

func (s commentService) GetCommentsByPostID(postID int64) ([]domain.Comment, error) {
	return s.repo.GetCommentsByPostID(postID)
}
