package requests

import (
	"trainee/internal/domain"
)

type LoginAuth struct {
	Email    string `json:"email" validate:"required,email" example:"example@email.com"`
	Password string `json:"password" validate:"required,gte=8" example:"01234567890"`
}

type RegisterAuth struct {
	Email    string `json:"email" validate:"required,email" example:"example@email.com"`
	Password string `json:"password" validate:"required,gte=8" example:"01234567890"`
	Name     string `json:"name" validate:"required,gte=3"`
}

func (r RegisterAuth) RegisterToUser() domain.User {
	return domain.User{
		Email:    r.Email,
		Name:     r.Name,
		Password: r.Password,
	}
}

type RegisterOauth2 struct {
	ID       string `json:"id" validate:"required"`
	Email    string `json:"email" validate:"required,email" example:"example@email.com"`
	Password string `json:"password" validate:"required,gte=8" example:"01234567890"`
	Name     string `json:"name" validate:"required,gte=3"`
}

func (r RegisterOauth2) RegisterToUser() domain.User {
	return domain.User{
		Email:    r.Email,
		Name:     r.Name,
		Password: r.Password,
	}
}
