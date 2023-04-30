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
	Migrate()

	Authenticate(email string, password string) (*model.Worker, error)

	GetWorkers() ([]*model.Worker, error)
	GetWorkerById(id model.WorkerID) (*model.Worker, error)
	CreateWorker(worker *model.Worker) error
	UpdateWorker(worker *model.Worker) error
	DeleteWorkerById(id model.WorkerID) error

	GetShifts(date *time.Time, span TimeSpan, workerId *model.WorkerID) ([]*model.Shift, error)
	GetShiftById(id model.ShiftID) (*model.Shift, error)
	CreateShift(shift *model.Shift) error
	UpdateShift(shift *model.Shift) error
	DeleteShiftById(id model.ShiftID) error

	CreateShiftAssignment(workerId model.WorkerID, shiftId model.ShiftID) error
	DeleteShiftAssignment(workerId model.WorkerID, shiftId model.ShiftID) error
}

func getSpanRange(date *time.Time, span TimeSpan) (time.Time, time.Time) {
	var intStart, intEnd time.Time
	y, m, d := date.Date()
	intStart = time.Date(y, m, d, 0, 0, 0, 0, nil)
	if span == DaySpan {
		intEnd = intStart.AddDate(0, 0, 1)
	} else {
		wd := intStart.Weekday()
		delta := (int(wd) - 1) % 7
		intStart = intStart.AddDate(0, 0, -delta)
		intEnd = intStart.AddDate(0, 0, 7)
	}
	return intStart, intEnd
}
