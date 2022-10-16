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
		val:     validators.NewCoordinateValidator(),
	}
}

func (c CommentHandler) SaveComment() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case "POST":
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
		default:
			writer.WriteHeader(http.StatusNotFound)
			return
		}
	}
}

func (c CommentHandler) GetComment() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case "GET":
			id, err := strconv.ParseInt(request.URL.Query().Get("id"), 10, 64)
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
		default:
			writer.WriteHeader(http.StatusNotFound)
			return
		}
	}
}

func (c CommentHandler) UpdateComment() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case "PUT":
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
		default:
			writer.WriteHeader(http.StatusNotFound)
			return
		}
	}
}

func (c CommentHandler) DeleteComment() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case "DELETE":
			id, err := strconv.ParseInt(request.URL.Query().Get("id"), 10, 64)
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
		default:
			writer.WriteHeader(http.StatusNotFound)
			return
		}
	}
}

//func (s service) DeleteComment(id int64) error {
//	//TODO implement me
//	panic("implement me")
//}
