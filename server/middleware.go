package server

import (
	"fmt"

	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"skybluetrades.net/work-planning-demo/api"
)

func CreateMiddleware(cfg *Config) ([]echo.MiddlewareFunc, error) {
	spec, err := api.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("Error loading swagger spec: %w", err)
	}
	spec.Servers = nil

	validator := middleware.OapiRequestValidatorWithOptions(spec,
		&middleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: NewAuthenticator(cfg),
			},
		})

	return []echo.MiddlewareFunc{validator}, nil
}
