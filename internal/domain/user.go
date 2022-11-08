package domain

import (
	"time"
	"trainee/internal/infra/http/response"
)

type User struct {
	ID          int64
	Email       string
	Name        string
	Password    string
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

func (u User) DomainToResponse() response.UserResponse {
	return response.UserResponse{
		ID:    u.ID,
		Email: u.Email,
		Name:  u.Name,
	}
}
