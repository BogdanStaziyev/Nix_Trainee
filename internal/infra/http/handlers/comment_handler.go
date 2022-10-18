package handlers

import (
	"fmt"
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
		val:     validators.NewCommentValidator(),
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
			BadRequest(writer, fmt.Errorf("expected 'id' to be an integer, was given: %s", chi.URLParam(request, "id")))
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

func (c CommentHandler) UpdateComment() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		comment, err := c.val.ValidateAndMap(request)
		if err != nil {
			log.Print(err)
			BadRequest(writer, err)
			return
		}
		comment, err = c.service.UpdateComment(comment)
		if err != nil {
			log.Print(err)
			InternalServerError(writer, err)
			return
		}
		success(writer, resources.MapDomainToCommentDTO(comment))
	}
}

func (c CommentHandler) DeleteComment() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(request, "id"), 10, 64)
		if err != nil {
			log.Print(err)
			BadRequest(writer, fmt.Errorf("expected 'id' to be an integer, was given: %s", chi.URLParam(request, "id")))
			return
		}
		err = c.service.DeleteComment(id)
		if err != nil {
			if strings.HasSuffix(err.Error(), "upper: no more rows in this result set") {
				NotFound(writer, err)
			} else {
				InternalServerError(writer, err)
			}
			return
		}
		ok(writer)
	}
}
