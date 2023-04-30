package server

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"skybluetrades.net/work-planning-demo/api"
	"skybluetrades.net/work-planning-demo/model"
	"skybluetrades.net/work-planning-demo/store"
)

// Delete an existing worker
// (DELETE /worker/{worker-id})
func (s *server) DeleteWorker(ctx echo.Context, workerId api.WorkerIdParam) error {
	err := s.db.DeleteWorkerById(model.WorkerID(workerId))
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

// Get all workers
// (GET /worker)
func (s *server) GetWorkers(ctx echo.Context) error {
	workers, err := s.db.GetWorkers()
	if err != nil {
		return err
	}

	// Convert the Worker models from the store into OpenAPI Worker
	// schema objects for return.
	ws := make([]*api.Worker, len(workers))
	for i, w := range workers {
		ws[i] = model.WorkerToAPI(w)
	}
	return ctx.JSON(http.StatusOK, ws)
}

// Create new worker
// (POST /worker)
func (s *server) CreateWorker(ctx echo.Context) error {
	var w api.Worker
	err := ctx.Bind(&w)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, "Invalid format for worker")
	}
	err = checkWorker(ctx, &w, false)
	if err != nil {
		return err
	}

	worker := model.WorkerFromAPI(&w)
	err = s.db.CreateWorker(worker)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, model.WorkerToAPI(worker))
}

// Update an existing worker
// (PUT /worker)
func (s *server) UpdateWorker(ctx echo.Context) error {
	var w api.Worker
	err := ctx.Bind(&w)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, "Invalid format for worker")
	}
	err = checkWorker(ctx, &w, true)
	if err != nil {
		return err
	}
	_, err = s.db.GetWorkerById(model.WorkerID(*w.Id))
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, "Unknown worker ID")
	}

	worker := model.WorkerFromAPI(&w)
	err = s.db.UpdateWorker(worker)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, model.WorkerToAPI(worker))
}

// Get a single worker
// (GET /worker/{worker-id})
func (s *server) GetWorker(ctx echo.Context, workerId api.WorkerIdParam) error {
	worker, err := s.db.GetWorkerById(model.WorkerID(workerId))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, model.WorkerToAPI(worker))
}

// Get schedule for a single worker
// (GET /worker/{worker-id}/schedule)
func (s *server) GetWorkerSchedule(ctx echo.Context,
	workerId api.WorkerIdParam, params api.GetWorkerScheduleParams) error {
	worker, err := s.db.GetWorkerById(model.WorkerID(workerId))
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
