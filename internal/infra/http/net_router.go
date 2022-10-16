package http

import (
	"net/http"
	"trainee/config/container"
)

func NetRouter(cont container.Container) http.Handler {

	router := http.NewServeMux()
	router.HandleFunc("/api/v1/comments/save", cont.CommentHandler.SaveComment())
	router.HandleFunc("/api/v1/comments/comment", cont.CommentHandler.GetComment())
	router.HandleFunc("/api/v1/comments/update", cont.CommentHandler.UpdateComment())
	router.HandleFunc("/api/v1/comments/delete", cont.CommentHandler.DeleteComment())

	return router
}
