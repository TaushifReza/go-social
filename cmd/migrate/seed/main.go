package main

import (
	"log"

	"github.com/TaushifReza/go-social/internal/db"
	"github.com/TaushifReza/go-social/internal/env"
	"github.com/TaushifReza/go-social/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://postgres:root786@172.21.64.1/social?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)

	db.Seed(store, conn)
}
