package server

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"skybluetrades.net/work-planning-demo/api"
	"skybluetrades.net/work-planning-demo/store"
)

type server struct {
	config *Config
	db     store.Store
}

func NewServer(cfg *Config, db store.Store, staticFiles *embed.FS) *echo.Echo {
	// Retrieve API spec.
	spec, err := api.GetSwagger()
	if err != nil {
		log.Fatalln("Error retrieving API spec: ", err)
	}

	// Authentication/validation middleware.
	mw, err := CreateMiddleware(spec, cfg)
	if err != nil {
		log.Fatalln("Error creating middleware: ", err)
	}

	// Set up Echo.
	e := echo.New()
	e.Debug = true
	e.Use(middleware.Logger())
	e.Use(mw...)

	// Register our server.
	srv := &server{
		config: cfg,
		db:     db,
	}
	api.RegisterHandlers(e, srv)

	// Routes for API docs.
	if staticFiles != nil {
		e.GET("/openapi3.json", func(ctx echo.Context) error {
			return ctx.JSON(http.StatusOK, spec)
		})
		fsys, err := fs.Sub(staticFiles, "static")
		if err != nil {
			log.Fatalln("Error setting up static files:", err)
		}
		e.GET("/*", echo.WrapHandler(http.FileServer(http.FS(fsys))))
	}

	// Return Echo instance to main!
	return e
}
