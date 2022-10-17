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

type PostHandler struct {
	service app.PostService
	val     validators.PostValidator
}

func NewPostHandler(s app.PostService) PostHandler {
	return PostHandler{
		service: s,
		val:     validators.NewPostValidator(),
	}
}

func (c PostHandler) SavePost() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case "POST":
			posts, err := c.val.ValidateAndMap(request)
			if err != nil {
				log.Print(err)
				BadRequest(writer, err)
				return
			}
			posts, err = c.service.SavePost(posts)
			if err != nil {
				log.Print(err)
				InternalServerError(writer, err)
				return
			}
			created(writer, resources.MapDomainToPostDTO(posts))
		default:
			writer.WriteHeader(http.StatusNotFound)
			return
		}
	}
}

func (c PostHandler) GetPost() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case "GET":
			id, err := strconv.ParseInt(request.URL.Query().Get("id"), 10, 64)
			if err != nil {
				log.Print(err)
				BadRequest(writer, fmt.Errorf("expected 'id' to be an integer, was given: %s", chi.URLParam(request, "id")))
				return
			}
			posts, err := c.service.GetPost(id)
			if err != nil {
				log.Print("PostService error", err)
				if strings.HasSuffix(err.Error(), "upper: no more rows in this result set") {
					NotFound(writer, err)
				} else {
					InternalServerError(writer, err)
				}
				return
			}
			success(writer, resources.MapDomainToPostDTO(posts))
			if err != nil {
				log.Print(err)
			}
		default:
			writer.WriteHeader(http.StatusNotFound)
			return
		}
	}
}

func (c PostHandler) UpdatePost() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case "PUT":
			post, err := c.val.ValidateAndMap(request)
			if err != nil {
				log.Print(err)
				BadRequest(writer, err)
				return
			}
			post, err = c.service.UpdatePost(post)
			if err != nil {
				log.Print(err)
				InternalServerError(writer, err)
				return
			}
			success(writer, resources.MapDomainToPostDTO(post))
		default:
			writer.WriteHeader(http.StatusNotFound)
			return
		}
	}
}

func (c PostHandler) DeletePost() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case "DELETE":
			id, err := strconv.ParseInt(request.URL.Query().Get("id"), 10, 64)
			if err != nil {
				log.Print(err)
				BadRequest(writer, fmt.Errorf("expected 'id' to be an integer, was given: %s", chi.URLParam(request, "id")))
				return
			}
			err = c.service.DeletePost(id)
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
