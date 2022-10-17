package validators

import "trainee/internal/domain"

type postRequest struct {
	UserId int    `json:"user_id" validate:"required"`
	Id     int    `json:"id"`
	Title  string `json:"title" validate:"required"`
	Body   string `json:"body" validate:"required"`
}

func mapPostRequestDomain(r postRequest) domain.Post {
	return domain.Post{
		UserId: r.UserId,
		Id:     r.Id,
		Title:  r.Title,
		Body:   r.Body,
	}
}
