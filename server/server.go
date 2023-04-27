package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"skybluetrades.net/work-planning-demo/api"
	"skybluetrades.net/work-planning-demo/store"
)

type server struct {
	config  *Config
	db      store.Store
	authKey string
}

func NewServer(cfg *Config) *server {
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

	if cfg.AddTestData {
		addTestData(db)
	}

	return &server{
		config:  cfg,
		db:      db,
		authKey: cfg.AuthKey,
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
	claims, err := ValidateRefreshToken(refresh.RefreshToken, s.authKey)
	if err != nil {
		return sendError(ctx, http.StatusForbidden, "Failed to refresh access token")
	}

	// Do a database lookup for the worker.
	worker, err := s.db.GetWorkerById(claims.ID)
	if err != nil {
		return sendError(ctx, http.StatusForbidden, "Failed to refresh access token (unknown user)")
	}

	// Generate new tokens.
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

// Get information about current user
// (GET /me)
func (s *server) GetMe(ctx echo.Context) error {
	claims := ctx.Get("claims").(*JWTClaim)
	worker, err := s.db.GetWorkerById(claims.ID)
	if err != nil {
		return sendError(ctx, http.StatusNotFound, "Worker record not found")
	}

	id := int64(worker.ID)
	w := api.Worker{
		Id:      &id,
		Email:   worker.Email,
		Name:    worker.Name,
		IsAdmin: worker.IsAdmin,
	}
	return ctx.JSON(http.StatusOK, w)
}

// Get schedule information for current user
// (GET /schedule)
func (s *server) GetSchedule(ctx echo.Context, params api.GetScheduleParams) error {
	return errors.New("NYI")
}

// Get shifts for a span of time
// (GET /shift)
func (s *server) GetShifts(ctx echo.Context, params api.GetShiftsParams) error {
	return errors.New("NYI")
}

// Create new shift
// (POST /shift)
func (s *server) CreateShift(ctx echo.Context) error { return errors.New("NYI") }

// Update an existing shift
// (PUT /shift)
func (s *server) UpdateShift(ctx echo.Context) error { return errors.New("NYI") }

// Delete an existing shift
// (DELETE /shift/{shift-id})
func (s *server) DeleteShift(ctx echo.Context, shiftId api.ShiftIdParam) error {
	return errors.New("NYI")
}

// Get a single shift
// (GET /shift/{shift-id})
func (s *server) GetShift(ctx echo.Context, shiftId api.ShiftIdParam) error {
	return errors.New("NYI")
}

// Delete an existing shift assignment
// (DELETE /shift/{shift-id}/assignment)
func (s *server) DeleteShiftAssignment(ctx echo.Context, shiftId api.ShiftIdParam) error {
	return errors.New("NYI")
}

// Create new shift assignment
// (POST /shift/{shift-id}/assignment)
func (s *server) CreateShiftAssignment(ctx echo.Context, shiftId api.ShiftIdParam) error {
	return errors.New("NYI")
}

// Delete an existing worker
// (DELETE /worker)
func (s *server) DeleteWorker(ctx echo.Context) error { return errors.New("NYI") }

// Get all workers
// (GET /worker)
func (s *server) GetWorkers(ctx echo.Context) error {
	workers, err := s.db.GetWorkers()
	if err != nil {
		return err
	}

	ws := make([]*api.Worker, len(workers))
	for i, w := range workers {
		id := int64(w.ID)
		rw := api.Worker{
			Id:      &id,
			Email:   w.Email,
			Name:    w.Name,
			IsAdmin: w.IsAdmin,
		}
		ws[i] = &rw
	}
	return ctx.JSON(http.StatusOK, ws)
}

// Create new worker
// (POST /worker)
func (s *server) CreateWorker(ctx echo.Context) error { return errors.New("NYI") }

// Update an existing worker
// (PUT /worker)
func (s *server) UpdateWorker(ctx echo.Context) error { return errors.New("NYI") }

// Get a single worker
// (GET /worker/{worker-id})
func (s *server) GetWorker(ctx echo.Context, workerId api.WorkerIdParam) error {
	return errors.New("NYI")
}
