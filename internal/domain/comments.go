package domain

import (
	"time"
	"trainee/internal/infra/http/response"
)

type Comment struct {
	ID          int64
	PostID      int64
	Name        string
	Email       string
	Body        string
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

func (c Comment) DomainToResponse() response.CommentResponse {
	return response.CommentResponse{
		ID:     c.ID,
		PostID: c.PostID,
		Name:   c.Name,
		Email:  c.Email,
		Body:   c.Body,
	}
}
