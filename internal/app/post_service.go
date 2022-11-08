package app

import (
	"github.com/golang-jwt/jwt"
	"log"
	"trainee/internal/domain"
	"trainee/internal/infra/database"
	"trainee/internal/infra/http/requests"
)

//go:generate mockery --dir . --name PostService --output ./mocks
type PostService interface {
	SavePost(postRequest requests.PostRequest, token *jwt.Token) (domain.Post, error)
	GetPost(id int64) (domain.Post, error)
	UpdatePost(postRequest requests.PostRequest, postID int64, token *jwt.Token) (domain.Post, error)
	DeletePost(id int64) error
	GetPostsByUser(userID int64) ([]domain.Post, error)
}

type postService struct {
	repo database.PostRepo
}

func NewPostService(repo database.PostRepo) PostService {
	return postService{
		repo: repo,
	}
}

func (s postService) SavePost(postRequest requests.PostRequest, token *jwt.Token) (domain.Post, error) {
	claim := token.Claims.(*JwtAccessClaim)
	userID := claim.ID
	domainPost := domain.Post{
		Title:  postRequest.Title,
		Body:   postRequest.Body,
		UserID: userID,
	}
	post, err := s.repo.SavePost(domainPost)
	if err != nil {
		log.Println(err)
		return domain.Post{}, err
	}
	return post, err
}

func (s postService) GetPost(id int64) (domain.Post, error) {
	return s.repo.GetPost(id)
}

func (s postService) UpdatePost(postRequest requests.PostRequest, postID int64, token *jwt.Token) (domain.Post, error) {
	claims := token.Claims.(*JwtAccessClaim)
	userID := claims.ID
	domainPost := domain.Post{
		Title:  postRequest.Title,
		Body:   postRequest.Body,
		UserID: userID,
		ID:     postID,
	}
	post, err := s.repo.UpdatePost(domainPost)
	if err != nil {
		log.Println(err)
		return domain.Post{}, err
	}
	return post, nil
}

func (s postService) DeletePost(id int64) error {
	return s.repo.DeletePost(id)
}

func (s postService) GetPostsByUser(userID int64) ([]domain.Post, error) {
	return s.repo.GetPostsByUser(userID)
}
