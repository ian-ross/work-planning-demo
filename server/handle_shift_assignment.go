package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"skybluetrades.net/work-planning-demo/api"
	"skybluetrades.net/work-planning-demo/model"
)

// Delete an existing shift assignment
// (DELETE /shift/{shift-id}/assignment)
func (s *server) DeleteShiftAssignment(ctx echo.Context, shiftId api.ShiftIdParam) error {
	worker, err := s.currentWorker(ctx)
	if err != nil {
		return err
	}

	err = s.db.DeleteShiftAssignment(worker.ID, model.ShiftID(shiftId))
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, "failed to delete assignment")
	}

	return ctx.NoContent(http.StatusNoContent)
}

// Create new shift assignment
// (POST /shift/{shift-id}/assignment)
func (s *server) CreateShiftAssignment(ctx echo.Context, shiftId api.ShiftIdParam) error {
	worker, err := s.currentWorker(ctx)
	if err != nil {
		return err
	}

	err = s.db.CreateShiftAssignment(worker.ID, model.ShiftID(shiftId))
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, "failed to delete assignment")
	}

	return ctx.NoContent(http.StatusNoContent)
}
