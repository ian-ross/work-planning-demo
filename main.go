package main

import (
	"fmt"
	"log"

	"github.com/dotenv-org/godotenvvault"
	"github.com/joeshaw/envdecode"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"skybluetrades.net/work-planning-demo/api"
	"skybluetrades.net/work-planning-demo/server"
)

func main() {
	// Load server config from environment, reading settings from .env
	// (or encrypted .env.vault if DOTENV_KEY is set).
	godotenvvault.Load()
	cfg := server.Config{}
	err := envdecode.StrictDecode(&cfg)

	// Authentication/validation middleware.
	mw, err := server.CreateMiddleware(&cfg)
	if err != nil {
		log.Fatalln("Error creating middleware: ", err)
	}

	// Set up Echo.
	e := echo.New()
	e.Debug = true
	e.Use(echomiddleware.Logger())
	e.Use(mw...)

	// Register our server.
	srv := server.NewServer(&cfg)
	api.RegisterHandlers(e, srv)

	// Off we go...
	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%d", cfg.Port)))
}
