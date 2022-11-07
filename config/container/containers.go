package container

import (
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"trainee/config"
	"trainee/internal/app"
	"trainee/internal/infra/database"
	"trainee/internal/infra/http/handlers"
)

type Container struct {
	Services
	Handlers
}

type Services struct {
	app.CommentService
	app.PostService
	app.UserService
	app.AuthService
}

type Handlers struct {
	handlers.CommentHandler
	handlers.PostHandler
	handlers.RegisterHandler
	handlers.OauthHandler
}

func New(conf config.Configuration) Container {
	sess := getDbSess(conf)

	userRepository := database.NewUSerRepo(sess)
	passwordGenerator := app.NewGeneratePasswordHash(bcrypt.DefaultCost)
	userService := app.NewUserService(userRepository, passwordGenerator)
	authService := app.NewAuthService(userService, conf)
	registerController := handlers.NewRegisterHandler(userService, authService)
	oauthController := handlers.NewOauthHandler(userService, authService)

	postRepository := database.NewPostRepository(sess)
	postService := app.NewPostService(postRepository)
	postHandler := handlers.NewPostHandler(postService)

	commentRepository := database.NewCommentRepository(sess)
	commentService := app.NewCommentService(commentRepository)
	commentHandler := handlers.NewCommentHandler(commentService, userService)

	return Container{
		Services: Services{
			commentService,
			postService,
			userService,
			authService,
		},
		Handlers: Handlers{
			commentHandler,
			postHandler,
			registerController,
			oauthController,
		},
	}
}

func getDbSess(conf config.Configuration) db.Session {
	sess, err := postgresql.Open(
		postgresql.ConnectionURL{
			User:     conf.DatabaseUser,
			Host:     conf.DatabaseHost,
			Password: conf.DatabasePassword,
			Database: conf.DatabaseName,
		})
	if err != nil {
		log.Fatalf("Unable to create new DB session: %q\n", err)
	}
	return sess
}
