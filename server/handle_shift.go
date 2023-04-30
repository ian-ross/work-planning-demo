package server

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"skybluetrades.net/work-planning-demo/api"
	"skybluetrades.net/work-planning-demo/model"
	"skybluetrades.net/work-planning-demo/store"
)

// Get shifts for a span of time
// (GET /shift)
func (s *server) GetShifts(ctx echo.Context, params api.GetShiftsParams) error {
	// Date defaults to nothing (which means today), span defaults to week.
	var date *time.Time
	if params.Date != nil {
		date = &params.Date.Time
	}
	span := store.WeekSpan
	if params.Span != nil && *params.Span == "day" {
		span = store.DaySpan
	}
	shifts, err := s.db.GetShifts(date, span, nil)
	if err != nil {
		return err
	}

	// Convert the Shift models from the store into OpenAPI Shift
	// schema objects for return.
	ss := make([]*api.Shift, len(shifts))
	for i, s := range shifts {
		ss[i] = model.ShiftToAPI(s)
	}
	return ctx.JSON(http.StatusOK, ss)
}

// Create new shift
// (POST /shift)
func (s *server) CreateShift(ctx echo.Context) error {
	var sh api.Shift
	err := ctx.Bind(&sh)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, "Invalid format for shift")
	}
	err = checkShift(ctx, &sh, false)
	if err != nil {
		return err
	}

	shift := model.ShiftFromAPI(&sh)
	err = s.db.CreateShift(shift)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, model.ShiftToAPI(shift))
}

// Update an existing shift
// (PUT /shift)
func (s *server) UpdateShift(ctx echo.Context) error {
	var sh api.Shift
	err := ctx.Bind(&sh)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, "Invalid format for shift")
	}
	err = checkShift(ctx, &sh, true)
	if err != nil {
		return err
	}
	_, err = s.db.GetShiftById(model.ShiftID(*sh.Id))
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, "Unknown shift ID")
	}

	shift := model.ShiftFromAPI(&sh)
	err = s.db.UpdateShift(shift)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, model.ShiftToAPI(shift))
}

// Delete an existing shift
// (DELETE /shift/{shift-id})
func (s *server) DeleteShift(ctx echo.Context, shiftId api.ShiftIdParam) error {
	err := s.db.DeleteShiftById(model.ShiftID(shiftId))
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

// Get a single shift
// (GET /shift/{shift-id})
func (s *server) GetShift(ctx echo.Context, shiftId api.ShiftIdParam) error {
	shift, err := s.db.GetShiftById(model.ShiftID(shiftId))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, model.ShiftToAPI(shift))
}
