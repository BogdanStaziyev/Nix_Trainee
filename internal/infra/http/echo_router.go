package http

import (
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
	"trainee/config"
	"trainee/config/container"
	"trainee/internal/app"
	"trainee/internal/infra/http/validators"

	_ "trainee/docs"
)

func EchoRouter(s *Server, cont container.Container) {

	e := s.Echo
	e.Use(middleware.Logger())
	e.Validator = validators.NewValidator()

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.POST("/register", cont.RegisterHandler.Register)
	e.POST("/login", cont.RegisterHandler.Login)

	v1 := e.Group("/api/v1")
	v1.GET("", PingHandler)

	conf := middleware.JWTConfig{
		Claims:     &app.JwtAccessClaim{},
		SigningKey: []byte(config.GetConfiguration().AccessSecret),
	}

	commRouter := v1.Group("/comments/")
	postRouter := v1.Group("/posts/")

	commRouter.Use(middleware.JWTWithConfig(conf))
	postRouter.Use(middleware.JWTWithConfig(conf))

	commRouter.POST("save/:post_id", cont.CommentHandler.SaveComment)
	commRouter.GET("comment/:id", cont.CommentHandler.GetComment)
	commRouter.PUT("update/:id", cont.CommentHandler.UpdateComment)
	commRouter.DELETE("delete/:id", cont.CommentHandler.DeleteComment)

	postRouter.POST("save", cont.PostHandler.SavePost)
	postRouter.GET("post/:id", cont.PostHandler.GetPost)
	postRouter.PUT("update/:id", cont.PostHandler.UpdatePost)
	postRouter.DELETE("delete/:id", cont.PostHandler.DeletePost)
}
