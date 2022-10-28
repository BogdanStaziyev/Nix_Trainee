package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
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
	exitCode := 0
	ctx, cancel := context.WithCancel(context.Background())

	//Recover
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("The sistem panicked!: %v\n", r)
			fmt.Printf("Stazk trace form panic: %s\n", string(debug.Stack()))
			exitCode = 1
		}
		os.Exit(exitCode)
	}()

	//Signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-c
		fmt.Printf("Received signal '%s', stopping... \n", sig.String())
		cancel()
		fmt.Printf("sent cancel to all threads...")
	}()

	var conf = config.GetConfiguration()

	err := database.Migrate(conf)
	if err != nil {
		log.Fatalf("Unable to apply migrations: %q\n", err)
	}

	cont := container.New(conf)

	// HTTP Server
	err = http.Server(
		ctx,
		http.EchoRouter(cont),
	)

	if err != nil {
		fmt.Printf("http server error: %s", err)
		exitCode = 2
		return
	}
}
