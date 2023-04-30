package server

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"skybluetrades.net/work-planning-demo/api"
	"skybluetrades.net/work-planning-demo/model"
	"skybluetrades.net/work-planning-demo/store"
)

// Get information about current user
// (GET /me)
func (s *server) GetMe(ctx echo.Context) error {
	worker, err := s.currentWorker(ctx)
	if err != nil {
		return err
	}

	// Make a worker and return it.
	return ctx.JSON(http.StatusOK, model.WorkerToAPI(worker))
}

// Get schedule information for current user
// (GET /schedule)
func (s *server) GetMeSchedule(ctx echo.Context, params api.GetMeScheduleParams) error {
	worker, err := s.currentWorker(ctx)
	if err != nil {
		return err
	}

	// Date defaults to nothing (which means today), span defaults to
	// week.
	var date *time.Time
	if params.Date != nil {
		date = &params.Date.Time
	}
	span := store.WeekSpan
	if params.Span != nil && *params.Span == "day" {
		span = store.DaySpan
	}
	shifts, err := s.db.GetShifts(date, span, &worker.ID)
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
