package http

import (
	"github.com/gorilla/mux"
	"net/http"
	"trainee/internal/infra/http/handlers"
)

func Router(commentHandler *handlers.CommentHandler) http.Handler {
	router := mux.NewRouter()

	// Health
	router.Group(func(healthRouter chi.Router) {
		healthRouter.Use(middleware.RedirectSlashes)

		healthRouter.Route("/ping", func(healthRouter chi.Router) {
			healthRouter.Get("/", PingHandler())

			healthRouter.Handle("/*", NotFoundJSON())
		})
	})

	router.Group(func(apiRouter chi.Router) {
		apiRouter.Use(middleware.RedirectSlashes)

		apiRouter.Route("/grass", func(apiRouter chi.Router) {
			AddEventRoutes(&apiRouter, eventController)

			apiRouter.Handle("/*", NotFoundJSON())
		})
	})

	return router
}

func AddEventRoutes(router *chi.Router, eventController *controllers.EventController) {
	(*router).Route("/events", func(apiRouter chi.Router) {
		apiRouter.Get(
			"/FindAllUsers",
			eventController.FindAll(),
		)
		apiRouter.Get(
			"/FindByName/{name}",
			eventController.FindByName(),
		)

		apiRouter.Get(
			"/Create/{name}/{age}/{city}/{country}",
			eventController.CreateUser(),
		)
		apiRouter.Get(
			"/UpdateData/{id}/{name}/{age}/{city}/{country}",
			eventController.UpdateById(),
		)
	})
}
