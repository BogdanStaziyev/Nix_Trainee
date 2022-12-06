package container

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"trainee/config"
	"trainee/internal/app"
	"trainee/internal/infra/database"
	"trainee/internal/infra/http/handlers"
	"trainee/middleware"
)

type Container struct {
	Services
	Handlers
	Middleware
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

type Middleware struct {
	middleware.AuthMiddleware
}

func New(conf config.Configuration) Container {
	sess := getDbSess(conf)
	newRedis := getRedis(conf)

	userRepository := database.NewUSerRepo(sess)
	passwordGenerator := app.NewGeneratePasswordHash(bcrypt.DefaultCost)
	userService := app.NewUserService(userRepository, passwordGenerator)
	authService := app.NewAuthService(userService, conf, newRedis)
	registerController := handlers.NewRegisterHandler(authService)
	oauthController := handlers.NewOauthHandler(userService, authService)

	postRepository := database.NewPostRepository(sess)
	postService := app.NewPostService(postRepository)

	commentRepository := database.NewCommentRepository(sess)
	commentService := app.NewCommentService(commentRepository, userService, postService)
	commentHandler := handlers.NewCommentHandler(commentService)

	postHandler := handlers.NewPostHandler(postService, commentService)

	authMiddleware := middleware.NewMiddleware(authService, newRedis)

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
		Middleware: Middleware{
			authMiddleware,
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

func getRedis(conf config.Configuration) *redis.Client {
	addr := fmt.Sprintf("%s:%s", conf.RedisHost, conf.RedisPort)
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}
