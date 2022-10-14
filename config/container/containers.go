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
}

type Handlers struct {
	handlers.CommentHandler
}

func New(conf config.Configuration) Container {
	sess := getDbSess(conf)

	commentRepository := database.NewCommentRepository(sess)
	commentService := app.NewCommentService(commentRepository)
	commentHandler := handlers.NewCommentHandler(commentService)

	return Container{
		Services: Services{
			commentService,
		},
		Handlers: Handlers{
			commentHandler,
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
