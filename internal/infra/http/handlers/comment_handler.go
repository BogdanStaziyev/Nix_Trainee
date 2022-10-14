package handlers

import (
	"log"
	"net/http"
	"trainee/internal/app"
	"trainee/internal/domain"
	"trainee/internal/infra/http/resources"
	"trainee/internal/infra/http/validators"
)

type CommentHandler struct {
	service app.CommentService
}

func NewCommentHandler(s app.CommentService) CommentHandler {
	return CommentHandler{
		service: s,
	}
}

func (c CommentHandler) SaveComment() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		comments, err := validators.Bind(request, validators.CommentRequest{}, domain.Comment{})
		if err != nil {
			log.Print(err)
			BadRequest(writer, err)
			return
		}
		comments, err = c.service.SaveComment(comments)
		if err != nil {
			log.Print(err)
			InternalServerError(writer, err)
			return
		}
		var commentsDto resources.CommentDTO
		created(writer, commentsDto.MapDomainToCommentDTO(comments))
	}
}
