package app

import (
	"fmt"
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
	UpdatePost(postRequest requests.PostRequest, postID int64) (domain.Post, error)
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
		return domain.Post{}, fmt.Errorf("service error save post: %w", err)
	}
	return post, nil
}

func (s postService) GetPost(id int64) (domain.Post, error) {
	post, err := s.repo.GetPost(id)
	if err != nil {
		return domain.Post{}, fmt.Errorf("service error get post: %w", err)
	}
	return post, nil
}

func (s postService) UpdatePost(postRequest requests.PostRequest, postID int64) (domain.Post, error) {
	post, err := s.repo.GetPost(postID)
	if err != nil {
		return domain.Post{}, fmt.Errorf("service error update post: %w", err)
	}

	post.Body = postRequest.Body
	post.Title = postRequest.Title

	post, err = s.repo.UpdatePost(post)
	if err != nil {
		log.Println(err)
		return domain.Post{}, fmt.Errorf("service error update post: %w", err)
	}
	return post, nil
}

func (s postService) DeletePost(id int64) error {
	_, err := s.repo.GetPost(id)
	if err != nil {
		return fmt.Errorf("service error delete post: %w", err)
	}
	err = s.repo.DeletePost(id)
	if err != nil {
		return fmt.Errorf("service error delete post: %w", err)
	}
	return err
}

func (s postService) GetPostsByUser(userID int64) ([]domain.Post, error) {
	posts, err := s.repo.GetPostsByUser(userID)
	if err != nil {
		return []domain.Post{}, fmt.Errorf("service error get posts by user id: %w", err)
	}
	return posts, nil
}
