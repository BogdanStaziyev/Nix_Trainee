package handlers

import (
	"net/http"
	"trainee/internal/app"
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

	}
}
