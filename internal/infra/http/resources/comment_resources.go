package resources

import "trainee/internal/domain"

type CommentDTO struct {
	Id     int64  `json:"id"`
	PostID int64  `json:"post_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

func (d CommentDTO) MapDomainToCommentDTO(comment domain.Comment) CommentDTO {
	d.Id = comment.Id
	d.PostID = comment.PostId
	d.Name = comment.Name
	d.Email = comment.Email
	d.Body = comment.Body
	return d
}
