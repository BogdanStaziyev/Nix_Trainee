package validators

import "trainee/internal/domain"

type commentRequest struct {
	ID     int64  `json:"id"`
	PostID int64  `json:"post_id" validate:"required"`
	Name   string `json:"name" validate:"required"`
	Email  string `json:"email" validate:"required"`
	Body   string `json:"body" validate:"required"`
}

func mapCommentRequestDomain(r commentRequest) domain.Comment {
	return domain.Comment{
		Id:     r.ID,
		PostId: r.PostID,
		Name:   r.Name,
		Email:  r.Email,
		Body:   r.Body,
	}
}
