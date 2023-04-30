package server

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"skybluetrades.net/work-planning-demo/api"
	"skybluetrades.net/work-planning-demo/model"
	"skybluetrades.net/work-planning-demo/store"
)

type server struct {
	config *Config
	db     store.Store
}

func NewServer(cfg *Config) *server {
	// Create a store for the server: options are a simple in-memory
	// store for testing, or Postgres (not implemented yet) determined
	// by the STORE_URL environment variable.
	var db store.Store
	var err error
	if cfg.StoreURL == "memory" {
		db, err = store.NewMemoryStore()
	} else {
		db, err = store.NewPostgresStore(cfg.StoreURL)
	}
	if err != nil {
		log.Fatalln("Error connecting to store: ", err)
	}
	db.Migrate()

	return &server{
		config: cfg,
		db:     db,
	}
}

// This function wraps sending of an error in the Error format, and
// handling the failure to marshal that.
func sendError(ctx echo.Context, code int, message string) error {
	errMsg := api.Error{
		Message: message,
	}
	err := ctx.JSON(code, errMsg)
	return err
}

// (POST /auth/login)
func (s *server) PostLogin(ctx echo.Context) error {
	var login api.Login
	err := ctx.Bind(&login)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, "Invalid format for login")
	}

	// Authenticate user.
	worker, err := s.db.Authenticate(login.Email, login.Password)
	if err != nil {
		return err
	}
	if worker == nil {
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
	var refresh api.CredentialsRefresh
	err := ctx.Bind(&refresh)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, "Invalid format for for token refresh")
	}

	// Decode and validate the refresh token from the request.
	claims, err := ValidateRefreshToken(refresh.RefreshToken, s.config.AuthKey)
	if err != nil {
		return sendError(ctx, http.StatusForbidden, "Failed to refresh access token")
	}

	// Do a database lookup for the worker.
	worker, err := s.db.GetWorkerById(claims.ID)
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
