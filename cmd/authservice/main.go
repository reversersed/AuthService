package main

import (
	"log"

	"github.com/reversersed/AuthService/internal/app"
)

func main() {

	app, err := app.New()
	if err != nil {
		log.Fatalf("error setting up the application: %v", err)
		return
	}
	if err := app.Run(); err != nil {
		log.Fatalf("error running the application: %v", err)
		return
	}
}
