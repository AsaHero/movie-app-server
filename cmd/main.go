package main

import (
	"github.com/AsaHero/movie-app-server/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables and start the application
	godotenv.Load()
	app.Run()
}
