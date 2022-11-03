package domain

import (
	"time"
	"trainee/internal/infra/http/response"
)

type Post struct {
	UserID      int64
	ID          int64
	Title       string
	Body        string
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

func (p Post) DomainToResponse() response.PostResponse {
	return response.PostResponse{
		ID:     p.ID,
		UserID: p.UserID,
		Title:  p.Title,
		Body:   p.Body,
	}
}
