package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
	"trainee/config/container"
	"trainee/internal/infra/http/handlers"
)

func Router(cont container.Container) http.Handler {

	router := chi.NewRouter()

	// Health
	router.Group(func(healthRouter chi.Router) {
		healthRouter.Use(middleware.RedirectSlashes)

		healthRouter.Route("/api/ping", func(healthRouter chi.Router) {
			healthRouter.Get("/", PingHandler())

			healthRouter.Handle("/*", NotFoundJSON())
		})
	})

	router.Group(func(apiRouter chi.Router) {
		apiRouter.Use(middleware.RedirectSlashes, cors.Handler(cors.Options{
			AllowedOrigins: []string{
				"https://*",
				"http://*",
			},
			AllowedMethods: []string{
				"GET",
				"POST",
				"PUT",
				"DELETE",
			},
			AllowedHeaders: []string{
				"Accept",
				"Content-Type",
			},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300,
		}))

		apiRouter.Route("/api/v1", func(apiRouter chi.Router) {
			apiRouter.Group(func(apiRouter chi.Router) {
				CommentRouter(apiRouter, cont.CommentHandler)
				apiRouter.Handle("/*", NotFoundJSON())

			})
			apiRouter.Handle("/*", NotFoundJSON())
		})
	})

	return router
}

func CommentRouter(router chi.Router, ch handlers.CommentHandler) {
	router.Route("/comments", func(apiRouter chi.Router) {
		apiRouter.Post(
			"/save",
			ch.SaveComment(),
		)
		apiRouter.Get(
			"/comment",
			ch.GetComment(),
		)
		apiRouter.Put(
			"/update",
			ch.UpdateComment(),
		)
		apiRouter.Delete(
			"/delete",
			ch.DeleteComment(),
		)
	})
}
