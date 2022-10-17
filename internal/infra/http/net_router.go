package http

import (
	"net/http"
	"trainee/config/container"
)

const (
	commentURL = "/api/v1/comments/"
	postURL    = "/api/v1/posts/"
)

func NetRouter(cont container.Container) http.Handler {

	router := http.NewServeMux()

	router.HandleFunc(commentURL+"save", cont.CommentHandler.SaveComment())
	router.HandleFunc(commentURL+"comment", cont.CommentHandler.GetComment())
	router.HandleFunc(commentURL+"update", cont.CommentHandler.UpdateComment())
	router.HandleFunc(commentURL+"delete", cont.CommentHandler.DeleteComment())

	router.HandleFunc(postURL+"save", cont.PostHandler.SavePost())
	router.HandleFunc(postURL+"post", cont.PostHandler.GetPost())
	router.HandleFunc(postURL+"update", cont.PostHandler.UpdatePost())
	router.HandleFunc(postURL+"delete", cont.PostHandler.DeletePost())

	return router
}
