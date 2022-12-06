package app

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"trainee/internal/domain"
	"trainee/internal/infra/database"
	"trainee/internal/infra/http/requests"
)

//go:generate mockery --dir . --name CommentService --output ./mocks
type CommentService interface {
	SaveComment(commentRequest requests.CommentRequest, postID int64, token *jwt.Token) (domain.Comment, error)
	GetComment(id int64) (domain.Comment, error)
	UpdateComment(commentRequest requests.CommentRequest, id int64) (domain.Comment, error)
	DeleteComment(id int64) error
	GetCommentsByPostID(postID int64, offset int) ([]domain.Comment, error)
}

type commentService struct {
	repo database.CommentRepo
	us   UserService
	ps   PostService
}

func NewCommentService(repo database.CommentRepo, us UserService, ps PostService) CommentService {
	return commentService{
		repo: repo,
		us:   us,
		ps:   ps,
	}
}

func (s commentService) SaveComment(commentRequest requests.CommentRequest, postID int64, token *jwt.Token) (domain.Comment, error) {
	claims := token.Claims.(*JwtTokenClaim)
	_, err := s.ps.GetPost(postID)
	if err != nil {
		return domain.Comment{}, fmt.Errorf("service error save comment: %w", err)
	}
	user, err := s.us.FindByID(claims.ID)
	if err != nil {
		return domain.Comment{}, fmt.Errorf("service error save comment: %w", err)
	}
	domainComment := domain.Comment{
		PostID: postID,
		Name:   user.Name,
		Email:  user.Email,
		Body:   commentRequest.Body,
	}
	comment, err := s.repo.SaveComment(domainComment)
	if err != nil {
		return domain.Comment{}, fmt.Errorf("service error save comment: %w", err)
	}
	return comment, nil
}

func (s commentService) GetComment(id int64) (domain.Comment, error) {
	comment, err := s.repo.GetComment(id)
	if err != nil {
		return domain.Comment{}, fmt.Errorf("service error get comment: %w", err)
	}
	return comment, nil
}

func (s commentService) UpdateComment(commentRequest requests.CommentRequest, id int64) (domain.Comment, error) {
	comment, err := s.repo.GetComment(id)
	if err != nil {
		return domain.Comment{}, fmt.Errorf("service error update comment: %w", err)
	}
	comment.Body = commentRequest.Body
	comment, err = s.repo.UpdateComment(comment)
	if err != nil {
		return domain.Comment{}, fmt.Errorf("service error update comment: %w", err)
	}
	return comment, nil
}

func (s commentService) DeleteComment(id int64) error {
	_, err := s.repo.GetComment(id)
	if err != nil {
		return fmt.Errorf("service error delete comment: %w", err)
	}
	err = s.repo.DeleteComment(id)
	if err != nil {
		return fmt.Errorf("service error delete comment: %w", err)
	}
	return nil
}

func (s commentService) GetCommentsByPostID(postID int64, offset int) ([]domain.Comment, error) {
	comments, err := s.repo.GetCommentsByPostID(postID, offset)
	if err != nil {
		return []domain.Comment{}, fmt.Errorf("service error get all comments by postID: %w", err)
	}
	return comments, nil
}
