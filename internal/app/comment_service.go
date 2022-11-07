package app

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
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
	GetCommentsByPostID(postID int64) ([]domain.Comment, error)
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
	claims := token.Claims.(*JwtAccessClaim)
	_, err := s.ps.GetPost(postID)
	if err != nil {
		log.Printf("SaveComment error, %s", err)
		return domain.Comment{}, fmt.Errorf("save comment error: %w", err)
	}
	user, err := s.us.FindByID(claims.ID)
	if err != nil {
		log.Printf("SaveComment error, %s", err)
		return domain.Comment{}, err
	}
	domainComment := domain.Comment{
		PostID: postID,
		Name:   user.Name,
		Email:  user.Email,
		Body:   commentRequest.Body,
	}
	return s.repo.SaveComment(domainComment)
}

func (s commentService) GetComment(id int64) (domain.Comment, error) {
	return s.repo.GetComment(id)
}

func (s commentService) UpdateComment(commentRequest requests.CommentRequest, id int64) (domain.Comment, error) {
	return s.repo.UpdateComment(commentRequest, id)
}

func (s commentService) DeleteComment(id int64) error {
	return s.repo.DeleteComment(id)
}

func (s commentService) GetCommentsByPostID(postID int64) ([]domain.Comment, error) {
	return s.repo.GetCommentsByPostID(postID)
}
