package main

import "github.com/reversersed/AuthService/internal/app"

func main() {

	app, err := app.New()
	if err != nil {
		return
	}
	if err := app.Run(); err != nil {
		return
	}
}
