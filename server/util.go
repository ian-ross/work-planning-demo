package server

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"skybluetrades.net/work-planning-demo/api"
	"skybluetrades.net/work-planning-demo/model"
)

// This function wraps sending of an error in the Error format, and
// handling the failure to marshal that.
func sendError(ctx echo.Context, code int, message string) error {
	errMsg := api.Error{
		Message: message,
	}
	err := ctx.JSON(code, errMsg)
	return err
}

func (s *server) currentWorker(ctx echo.Context) (*model.Worker, error) {
	// We need the idea of the "current user" independent of the
	// authentication flow. We can get the user ID associated with the
	// JWT that was used for authentication from the token claims that
	// we stored in the Echo context in the authentication middleware.
	claims := ctx.Get("claims").(*JWTClaim)
	worker, err := s.db.GetWorkerById(claims.ID)
	if err != nil {
		return nil, sendError(ctx, http.StatusNotFound, "Worker record not found")
	}
	return worker, nil
}

func checkWorker(ctx echo.Context, w *api.Worker, needId bool) error {
	if needId && w.Id == nil {
		return sendError(ctx, http.StatusBadRequest, "Missing worker ID")
	}
	if w.Password == nil || len(*w.Password) == 0 {
		return sendError(ctx, http.StatusBadRequest, "Missing password for worker")
	}
	if strings.TrimSpace(w.Name) == "" {
		return sendError(ctx, http.StatusBadRequest, "Missing name for worker")
	}
	if strings.TrimSpace(w.Email) == "" {
		return sendError(ctx, http.StatusBadRequest, "Missing email for worker")
	}
	return nil
}

func checkShift(ctx echo.Context, s *api.Shift, needId bool) error {
	if needId && s.Id == nil {
		return sendError(ctx, http.StatusBadRequest, "Missing shift ID")
	}
	if s.Capacity <= 0 {
		return sendError(ctx, http.StatusBadRequest, "Bad capacity value for shift")
	}
	if s.StartTime.After(s.EndTime) {
		return sendError(ctx, http.StatusBadRequest, "Bad time range for shift")
	}
	return nil
}
