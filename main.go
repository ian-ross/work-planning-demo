package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/dotenv-org/godotenvvault"
	"github.com/joeshaw/envdecode"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"skybluetrades.net/work-planning-demo/api"
	"skybluetrades.net/work-planning-demo/server"
)

//go:embed static
var staticFiles embed.FS

func main() {
	// Load server config from environment, reading settings from .env
	// (or encrypted .env.vault if DOTENV_KEY is set).
	godotenvvault.Load()
	cfg := server.Config{}
	err := envdecode.StrictDecode(&cfg)

	// Retrieve API spec.
	spec, err := api.GetSwagger()
	if err != nil {
		log.Fatalln("Error setting up API docs: ", err)
	}

	// Authentication/validation middleware.
	mw, err := server.CreateMiddleware(spec, &cfg)
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

	// Routes for API docs.
	e.GET("/openapi3.json", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, spec)
	})
	fsys, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalln("Error setting up static files:", err)
	}
	e.GET("/*", echo.WrapHandler(http.FileServer(http.FS(fsys))))

	// Off we go...
	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%d", cfg.Port)))
}
