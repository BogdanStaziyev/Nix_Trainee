package main

import (
	"log"
	"trainee/config"
	"trainee/config/container"
	"trainee/internal/infra/database"
	"trainee/internal/infra/http"
)

// @title 		NIX TRAINEE PROGRAM Demo App
// @version 	V1.echo
// @description REST service for NIX TRAINEE PROGRAM

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @host 		localhost:8080
// @BasePath 	/api/v1
func main() {
	var conf = config.GetConfiguration()

	err := database.Migrate(conf)
	if err != nil {
		log.Fatalf("Unable to apply migrations: %q\n", err)
	}

	cont := container.New(conf)

	// Echo Server
	srv := http.NewServer()

	http.EchoRouter(srv, cont)

	err = srv.Start()
	if err != nil {
		log.Fatal("Port already used")
	}
}
