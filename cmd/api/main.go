package main

import (
	"log"

	"github.com/TaushifReza/go-social/internal/env"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic("No .env file found (using environment variables)")
	}

	config := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	app := &application{
		config: config,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
