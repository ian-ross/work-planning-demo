package store

import (
	"time"

	"skybluetrades.net/work-planning-demo/model"
)

type TimeSpan int

const (
	WeekSpan TimeSpan = iota
	DaySpan  TimeSpan = iota
)

type Store interface {
	Authenticate(email string, password string) (*model.Worker, error)

	GetWorkers() ([]*model.Worker, error)
	GetWorkerById(id model.WorkerID) (*model.Worker, error)
	CreateWorker(worker *model.Worker) (*model.Worker, error)
	UpdateWorker(worker *model.Worker) (*model.Worker, error)
	DeleteWorker(worker *model.Worker) error

	GetShifts(date *time.Time, span TimeSpan) ([]*model.Shift, error)
	GetShiftById(id model.ShiftID) (*model.Shift, error)
	CreateShift(shift *model.Shift) (*model.Shift, error)
	UpdateShift(shift *model.Shift) (*model.Shift, error)
	DeleteShift(shift *model.Shift) error

	CreateShiftAssignment(workerId model.WorkerID, shiftId model.ShiftID) (*model.ShiftAssignment, error)
	DeleteShiftAssignment(assignment *model.ShiftAssignment) error
}
