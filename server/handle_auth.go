package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"skybluetrades.net/work-planning-demo/api"
)

// (POST /auth/login)
func (s *server) PostLogin(ctx echo.Context) error {
	var login api.Login
	err := ctx.Bind(&login)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, "Invalid format for login")
	}

	// Authenticate user.
	worker, err := s.db.Authenticate(login.Email, login.Password)
	if err != nil || worker == nil {
		return sendError(ctx, http.StatusForbidden, "Invalid login credentials")
	}

	// Generate JWTs and return credentials.
	accessToken, refreshToken, err := GenerateTokens(worker, s.config)
	if err != nil {
		return err
	}
	creds := api.Credentials{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return ctx.JSON(http.StatusOK, creds)
}

// (POST /auth/logout)
func (s *server) PostLogout(ctx echo.Context) error {
	// TODO: MAINTAIN LIST OF RETIRED ACCESS AND REFRESH TOKENS TO FORCE
	// RE-LOGIN? I NEVER KNOW WHAT TO DO WITH THIS...
	return ctx.NoContent(http.StatusNoContent)
}

// (POST /auth/refresh_token)
func (s *server) PostRefreshToken(ctx echo.Context) error {
	fmt.Println("===> 1")
	var refresh api.CredentialsRefresh
	err := ctx.Bind(&refresh)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, "Invalid format for for token refresh")
	}
	fmt.Println("===> 2: refresh = ", refresh)

	// Decode and validate the refresh token from the request.
	claims, err := ValidateRefreshToken(refresh.RefreshToken, s.config.AuthKey)
	fmt.Println("===> 3: claims = ", claims, "   err = ", err, "   claims.ID = ", claims.ID)
	if err != nil {
		return sendError(ctx, http.StatusForbidden, "Failed to refresh access token")
	}

	// Do a database lookup for the worker.
	worker, err := s.db.GetWorkerById(claims.ID)
	fmt.Println("===> 4: worker = ", worker, "   err = ", err)
	if err != nil {
		return sendError(ctx, http.StatusForbidden, "Failed to refresh access token (unknown user)")
	}

	// Generate and send new tokens.
	accessToken, refreshToken, err := GenerateTokens(worker, s.config)
	if err != nil {
		return err
	}
	creds := api.Credentials{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return ctx.JSON(http.StatusOK, creds)
}
