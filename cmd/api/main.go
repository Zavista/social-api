package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/zavista/social-api/internal/env"
	"github.com/zavista/social-api/internal/store"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	store := store.NewPostgresStorage(nil)

	app := &application{
		config: cfg,
		store:  store,
	}

	if err := app.run(); err != nil {
		log.Fatal(err)
	}

}
