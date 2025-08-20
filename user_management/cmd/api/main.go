package main

import (
	"enterdev.com.vn/user_management/internal/app"
	"enterdev.com.vn/user_management/internal/config"
)

func main() {
	// initialize configuration
	config := config.NewConfig()

	// init app
	app := app.NewApplication(config)

	// start server
	if err := app.Run(); err != nil {
		panic(err)
	}
}
