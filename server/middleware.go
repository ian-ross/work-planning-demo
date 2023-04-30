package server

import (
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
)

func CreateMiddleware(spec *openapi3.T, cfg *Config) ([]echo.MiddlewareFunc, error) {
	validator := middleware.OapiRequestValidatorWithOptions(spec,
		&middleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: NewAuthenticator(cfg),
			},
			Skipper: func(ctx echo.Context) bool {
				// Skip checks for static files.
				p := ctx.Path()
				return p == "/openapi3.json" || p == "/*"
			},
		})

	return []echo.MiddlewareFunc{validator}, nil
}
