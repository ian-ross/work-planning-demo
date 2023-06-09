// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for SpanLength.
const (
	SpanLengthDay  SpanLength = "day"
	SpanLengthWeek SpanLength = "week"
)

// Defines values for GetMeScheduleParamsSpan.
const (
	GetMeScheduleParamsSpanDay  GetMeScheduleParamsSpan = "day"
	GetMeScheduleParamsSpanWeek GetMeScheduleParamsSpan = "week"
)

// Defines values for GetShiftsParamsSpan.
const (
	GetShiftsParamsSpanDay  GetShiftsParamsSpan = "day"
	GetShiftsParamsSpanWeek GetShiftsParamsSpan = "week"
)

// Defines values for GetWorkerScheduleParamsSpan.
const (
	GetWorkerScheduleParamsSpanDay  GetWorkerScheduleParamsSpan = "day"
	GetWorkerScheduleParamsSpanWeek GetWorkerScheduleParamsSpan = "week"
)

// Credentials defines model for Credentials.
type Credentials struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// CredentialsRefresh defines model for CredentialsRefresh.
type CredentialsRefresh struct {
	RefreshToken string `json:"refresh_token"`
}

// Error defines model for Error.
type Error struct {
	Message string `json:"message"`
}

// Login defines model for Login.
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Shift defines model for Shift.
type Shift struct {
	AssignedWorkers *[]WorkerId `json:"assigned_workers,omitempty"`
	Capacity        int32       `json:"capacity"`
	EndTime         time.Time   `json:"end_time"`
	Id              *ShiftId    `json:"id,omitempty"`
	StartTime       time.Time   `json:"start_time"`
}

// ShiftId defines model for ShiftId.
type ShiftId = int64

// Worker defines model for Worker.
type Worker struct {
	Email    string    `json:"email"`
	Id       *WorkerId `json:"id,omitempty"`
	IsAdmin  bool      `json:"is_admin"`
	Name     string    `json:"name"`
	Password *string   `json:"password,omitempty"`
}

// WorkerId defines model for WorkerId.
type WorkerId = int64

// ShiftIdParam defines model for ShiftIdParam.
type ShiftIdParam = ShiftId

// SpanDate defines model for SpanDate.
type SpanDate = openapi_types.Date

// SpanLength defines model for SpanLength.
type SpanLength string

// WorkerIdParam defines model for WorkerIdParam.
type WorkerIdParam = WorkerId

// GetMeScheduleParams defines parameters for GetMeSchedule.
type GetMeScheduleParams struct {
	// Date Date including in weekly schedule to fetch (defaults to today)
	Date *SpanDate `form:"date,omitempty" json:"date,omitempty"`

	// Span Span of schedule ("week" or "day", defaults to "week")
	Span *GetMeScheduleParamsSpan `form:"span,omitempty" json:"span,omitempty"`
}

// GetMeScheduleParamsSpan defines parameters for GetMeSchedule.
type GetMeScheduleParamsSpan string

// GetShiftsParams defines parameters for GetShifts.
type GetShiftsParams struct {
	// Date Date including in weekly schedule to fetch (defaults to today)
	Date *SpanDate `form:"date,omitempty" json:"date,omitempty"`

	// Span Span of schedule ("week" or "day", defaults to "week")
	Span *GetShiftsParamsSpan `form:"span,omitempty" json:"span,omitempty"`
}

// GetShiftsParamsSpan defines parameters for GetShifts.
type GetShiftsParamsSpan string

// GetWorkerScheduleParams defines parameters for GetWorkerSchedule.
type GetWorkerScheduleParams struct {
	// Date Date including in weekly schedule to fetch (defaults to today)
	Date *SpanDate `form:"date,omitempty" json:"date,omitempty"`

	// Span Span of schedule ("week" or "day", defaults to "week")
	Span *GetWorkerScheduleParamsSpan `form:"span,omitempty" json:"span,omitempty"`
}

// GetWorkerScheduleParamsSpan defines parameters for GetWorkerSchedule.
type GetWorkerScheduleParamsSpan string

// PostLoginJSONRequestBody defines body for PostLogin for application/json ContentType.
type PostLoginJSONRequestBody = Login

// PostRefreshTokenJSONRequestBody defines body for PostRefreshToken for application/json ContentType.
type PostRefreshTokenJSONRequestBody = CredentialsRefresh

// CreateShiftJSONRequestBody defines body for CreateShift for application/json ContentType.
type CreateShiftJSONRequestBody = Shift

// UpdateShiftJSONRequestBody defines body for UpdateShift for application/json ContentType.
type UpdateShiftJSONRequestBody = Shift

// CreateWorkerJSONRequestBody defines body for CreateWorker for application/json ContentType.
type CreateWorkerJSONRequestBody = Worker

// UpdateWorkerJSONRequestBody defines body for UpdateWorker for application/json ContentType.
type UpdateWorkerJSONRequestBody = Worker

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /auth/login)
	PostLogin(ctx echo.Context) error

	// (POST /auth/logout)
	PostLogout(ctx echo.Context) error

	// (POST /auth/refresh_token)
	PostRefreshToken(ctx echo.Context) error
	// Get information about current user
	// (GET /me)
	GetMe(ctx echo.Context) error
	// Get schedule information for current user
	// (GET /me/schedule)
	GetMeSchedule(ctx echo.Context, params GetMeScheduleParams) error
	// Get shifts for a span of time
	// (GET /shift)
	GetShifts(ctx echo.Context, params GetShiftsParams) error
	// Create new shift
	// (POST /shift)
	CreateShift(ctx echo.Context) error
	// Update an existing shift
	// (PUT /shift)
	UpdateShift(ctx echo.Context) error
	// Delete an existing shift
	// (DELETE /shift/{shift-id})
	DeleteShift(ctx echo.Context, shiftId ShiftIdParam) error
	// Get a single shift
	// (GET /shift/{shift-id})
	GetShift(ctx echo.Context, shiftId ShiftIdParam) error
	// Delete an existing shift assignment
	// (DELETE /shift/{shift-id}/assignment)
	DeleteShiftAssignment(ctx echo.Context, shiftId ShiftIdParam) error
	// Create new shift assignment
	// (POST /shift/{shift-id}/assignment)
	CreateShiftAssignment(ctx echo.Context, shiftId ShiftIdParam) error
	// Get all workers
	// (GET /worker)
	GetWorkers(ctx echo.Context) error
	// Create new worker
	// (POST /worker)
	CreateWorker(ctx echo.Context) error
	// Update an existing worker
	// (PUT /worker)
	UpdateWorker(ctx echo.Context) error
	// Delete an existing worker
	// (DELETE /worker/{worker-id})
	DeleteWorker(ctx echo.Context, workerId WorkerIdParam) error
	// Get a single worker
	// (GET /worker/{worker-id})
	GetWorker(ctx echo.Context, workerId WorkerIdParam) error
	// Get schedule for a single worker
	// (GET /worker/{worker-id}/schedule)
	GetWorkerSchedule(ctx echo.Context, workerId WorkerIdParam, params GetWorkerScheduleParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// PostLogin converts echo context to params.
func (w *ServerInterfaceWrapper) PostLogin(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostLogin(ctx)
	return err
}

// PostLogout converts echo context to params.
func (w *ServerInterfaceWrapper) PostLogout(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostLogout(ctx)
	return err
}

// PostRefreshToken converts echo context to params.
func (w *ServerInterfaceWrapper) PostRefreshToken(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostRefreshToken(ctx)
	return err
}

// GetMe converts echo context to params.
func (w *ServerInterfaceWrapper) GetMe(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetMe(ctx)
	return err
}

// GetMeSchedule converts echo context to params.
func (w *ServerInterfaceWrapper) GetMeSchedule(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetMeScheduleParams
	// ------------- Optional query parameter "date" -------------

	err = runtime.BindQueryParameter("form", true, false, "date", ctx.QueryParams(), &params.Date)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter date: %s", err))
	}

	// ------------- Optional query parameter "span" -------------

	err = runtime.BindQueryParameter("form", true, false, "span", ctx.QueryParams(), &params.Span)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter span: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetMeSchedule(ctx, params)
	return err
}

// GetShifts converts echo context to params.
func (w *ServerInterfaceWrapper) GetShifts(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetShiftsParams
	// ------------- Optional query parameter "date" -------------

	err = runtime.BindQueryParameter("form", true, false, "date", ctx.QueryParams(), &params.Date)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter date: %s", err))
	}

	// ------------- Optional query parameter "span" -------------

	err = runtime.BindQueryParameter("form", true, false, "span", ctx.QueryParams(), &params.Span)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter span: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetShifts(ctx, params)
	return err
}

// CreateShift converts echo context to params.
func (w *ServerInterfaceWrapper) CreateShift(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{"admin"})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateShift(ctx)
	return err
}

// UpdateShift converts echo context to params.
func (w *ServerInterfaceWrapper) UpdateShift(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{"admin"})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.UpdateShift(ctx)
	return err
}

// DeleteShift converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteShift(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "shift-id" -------------
	var shiftId ShiftIdParam

	err = runtime.BindStyledParameterWithLocation("simple", false, "shift-id", runtime.ParamLocationPath, ctx.Param("shift-id"), &shiftId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter shift-id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{"admin"})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteShift(ctx, shiftId)
	return err
}

// GetShift converts echo context to params.
func (w *ServerInterfaceWrapper) GetShift(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "shift-id" -------------
	var shiftId ShiftIdParam

	err = runtime.BindStyledParameterWithLocation("simple", false, "shift-id", runtime.ParamLocationPath, ctx.Param("shift-id"), &shiftId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter shift-id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetShift(ctx, shiftId)
	return err
}

// DeleteShiftAssignment converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteShiftAssignment(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "shift-id" -------------
	var shiftId ShiftIdParam

	err = runtime.BindStyledParameterWithLocation("simple", false, "shift-id", runtime.ParamLocationPath, ctx.Param("shift-id"), &shiftId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter shift-id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteShiftAssignment(ctx, shiftId)
	return err
}

// CreateShiftAssignment converts echo context to params.
func (w *ServerInterfaceWrapper) CreateShiftAssignment(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "shift-id" -------------
	var shiftId ShiftIdParam

	err = runtime.BindStyledParameterWithLocation("simple", false, "shift-id", runtime.ParamLocationPath, ctx.Param("shift-id"), &shiftId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter shift-id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateShiftAssignment(ctx, shiftId)
	return err
}

// GetWorkers converts echo context to params.
func (w *ServerInterfaceWrapper) GetWorkers(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{"admin"})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetWorkers(ctx)
	return err
}

// CreateWorker converts echo context to params.
func (w *ServerInterfaceWrapper) CreateWorker(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{"admin"})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateWorker(ctx)
	return err
}

// UpdateWorker converts echo context to params.
func (w *ServerInterfaceWrapper) UpdateWorker(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{"admin"})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.UpdateWorker(ctx)
	return err
}

// DeleteWorker converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteWorker(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "worker-id" -------------
	var workerId WorkerIdParam

	err = runtime.BindStyledParameterWithLocation("simple", false, "worker-id", runtime.ParamLocationPath, ctx.Param("worker-id"), &workerId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter worker-id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{"admin"})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteWorker(ctx, workerId)
	return err
}

// GetWorker converts echo context to params.
func (w *ServerInterfaceWrapper) GetWorker(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "worker-id" -------------
	var workerId WorkerIdParam

	err = runtime.BindStyledParameterWithLocation("simple", false, "worker-id", runtime.ParamLocationPath, ctx.Param("worker-id"), &workerId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter worker-id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{"admin"})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetWorker(ctx, workerId)
	return err
}

// GetWorkerSchedule converts echo context to params.
func (w *ServerInterfaceWrapper) GetWorkerSchedule(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "worker-id" -------------
	var workerId WorkerIdParam

	err = runtime.BindStyledParameterWithLocation("simple", false, "worker-id", runtime.ParamLocationPath, ctx.Param("worker-id"), &workerId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter worker-id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{"admin"})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetWorkerScheduleParams
	// ------------- Optional query parameter "date" -------------

	err = runtime.BindQueryParameter("form", true, false, "date", ctx.QueryParams(), &params.Date)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter date: %s", err))
	}

	// ------------- Optional query parameter "span" -------------

	err = runtime.BindQueryParameter("form", true, false, "span", ctx.QueryParams(), &params.Span)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter span: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetWorkerSchedule(ctx, workerId, params)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/auth/login", wrapper.PostLogin)
	router.POST(baseURL+"/auth/logout", wrapper.PostLogout)
	router.POST(baseURL+"/auth/refresh_token", wrapper.PostRefreshToken)
	router.GET(baseURL+"/me", wrapper.GetMe)
	router.GET(baseURL+"/me/schedule", wrapper.GetMeSchedule)
	router.GET(baseURL+"/shift", wrapper.GetShifts)
	router.POST(baseURL+"/shift", wrapper.CreateShift)
	router.PUT(baseURL+"/shift", wrapper.UpdateShift)
	router.DELETE(baseURL+"/shift/:shift-id", wrapper.DeleteShift)
	router.GET(baseURL+"/shift/:shift-id", wrapper.GetShift)
	router.DELETE(baseURL+"/shift/:shift-id/assignment", wrapper.DeleteShiftAssignment)
	router.POST(baseURL+"/shift/:shift-id/assignment", wrapper.CreateShiftAssignment)
	router.GET(baseURL+"/worker", wrapper.GetWorkers)
	router.POST(baseURL+"/worker", wrapper.CreateWorker)
	router.PUT(baseURL+"/worker", wrapper.UpdateWorker)
	router.DELETE(baseURL+"/worker/:worker-id", wrapper.DeleteWorker)
	router.GET(baseURL+"/worker/:worker-id", wrapper.GetWorker)
	router.GET(baseURL+"/worker/:worker-id/schedule", wrapper.GetWorkerSchedule)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"2P7MIgAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
