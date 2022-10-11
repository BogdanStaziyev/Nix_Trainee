package main

import (
	"log"
	myHttp "trainee/internal/infra/http"
)

func main() {
	srv := new(myHttp.Server)
	if err := srv.Run("8080"); err != nil {
		log.Fatalf("error running http Server: %s", err.Error())
	}
}
