package resources

import "trainee/internal/domain"

type CommentDTO struct {
	Id     int64  `json:"id"`
	PostID int64  `json:"post_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

func MapDomainToCommentDTO(comment domain.Comment) CommentDTO {
	return CommentDTO{
		Id:     comment.Id,
		PostID: comment.PostId,
		Name:   comment.Name,
		Email:  comment.Email,
		Body:   comment.Body,
	}
}
