package validators

import "trainee/internal/domain"

type CommentRequest struct {
	PostID int64  `json:"post_id" validate:"required"`
	Name   string `json:"name" validate:"required"`
	Email  string `json:"email" validate:"required"`
	Body   string `json:"body" validate:"required"`
}

func (r CommentRequest) ToDomainModel() (interface{}, error) {
	return domain.Comment{
		PostId: r.PostID,
		Name:   r.Name,
		Email:  r.Email,
		Body:   r.Body,
	}, nil
}
