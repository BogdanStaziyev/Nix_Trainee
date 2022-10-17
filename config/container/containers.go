package container

import (
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
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
}

type Handlers struct {
	handlers.CommentHandler
	handlers.PostHandler
}

func New(conf config.Configuration) Container {
	sess := getDbSess(conf)

	commentRepository := database.NewCommentRepository(sess)
	commentService := app.NewCommentService(commentRepository)
	commentHandler := handlers.NewCommentHandler(commentService)

	postRepository := database.NewPostRepository(sess)
	postService := app.NewPostService(postRepository)
	postHandler := handlers.NewPostHandler(postService)

	return Container{
		Services: Services{
			commentService,
			postService,
		},
		Handlers: Handlers{
			commentHandler,
			postHandler,
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
