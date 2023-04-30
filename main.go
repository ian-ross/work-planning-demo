package main

import (
	"embed"
	"fmt"
	"log"

	"github.com/dotenv-org/godotenvvault"
	"github.com/joeshaw/envdecode"
	"skybluetrades.net/work-planning-demo/server"
	"skybluetrades.net/work-planning-demo/store"
)

//go:embed static
var staticFiles embed.FS

func main() {
	// Load server config from environment, reading settings from .env
	// (or encrypted .env.vault if DOTENV_KEY is set).
	godotenvvault.Load()
	cfg := server.Config{}
	err := envdecode.StrictDecode(&cfg)
	if err != nil {
		log.Fatalln("Error retrieving environment settings: ", err)
	}

	// Create a store for the server: options are a simple in-memory
	// store for testing, or Postgres (not implemented yet) determined
	// by the STORE_URL environment variable.
	var db store.Store
	if cfg.StoreURL == "memory" {
		db, err = store.NewMemoryStore()
	} else {
		db, err = store.NewPostgresStore(cfg.StoreURL)
	}
	if err != nil {
		log.Fatalln("Error connecting to store: ", err)
	}
	db.Migrate()

	// Set up Echo server.
	e := server.NewServer(&cfg, db, &staticFiles)

	// Off we go...
	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%d", cfg.Port)))
}
