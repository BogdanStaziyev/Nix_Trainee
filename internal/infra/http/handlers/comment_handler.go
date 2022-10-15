package handlers

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
	"strings"
	"trainee/internal/app"
	"trainee/internal/infra/http/resources"
	"trainee/internal/infra/http/validators"
)

type CommentHandler struct {
	service app.CommentService
	val     validators.CommentValidator
}

func NewCommentHandler(s app.CommentService) CommentHandler {
	return CommentHandler{
		service: s,
		val:     validators.NewCoordinateValidator(),
	}
}

func (c CommentHandler) SaveComment() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		comments, err := c.val.ValidateAndMap(request)
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
		created(writer, resources.MapDomainToCommentDTO(comments))
	}
}

func (c CommentHandler) GetComment() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(request, "id"), 10, 64)
		if err != nil {
			log.Print(err)
			BadRequest(writer, err)
			return
		}
		comments, err := c.service.GetComment(id)
		if err != nil {
			log.Print("commentService error", err)
			if strings.HasSuffix(err.Error(), "upper: no more rows in this result set") {
				NotFound(writer, err)
			} else {
				InternalServerError(writer, err)
			}
			return
		}
		success(writer, resources.MapDomainToCommentDTO(comments))
		if err != nil {
			log.Print(err)
		}
	}
}

//
//func (s service) GetComment(id int64) (domain.Comment, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s service) UpdateComment(comment domain.Comment) (domain.Comment, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s service) DeleteComment(id int64) error {
//	//TODO implement me
//	panic("implement me")
//}
