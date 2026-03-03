package main

import (
	"log"

	"github.com/TaushifReza/go-social/internal/db"
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
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://postgres:root786@172.21.64.1/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env:     env.GetString("ENV", "development"),
		version: env.GetString("VERSION", "0.0.1"),
	}

	db, err := db.New(
		config.db.addr,
		config.db.maxOpenConns,
		config.db.maxIdleConns,
		config.db.maxIdleTime,
	)

	if err != nil {
		log.Panic("database connection error: ", err)
	}

	store := store.NewStorage(db)

	app := &application{
		config: config,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
