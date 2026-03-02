package main

import (
	"log"

	"github.com/TaushifReza/go-social/internal/env"
	"github.com/TaushifReza/go-social/internal/store"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic("No .env file found (using environment variables)")
	}

	config := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	store := store.NewStorage(nil)

	app := &application{
		config: config,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
