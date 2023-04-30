// Code generated by mockery v2.26.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	model "skybluetrades.net/work-planning-demo/model"

	store "skybluetrades.net/work-planning-demo/store"

	time "time"
)

// Store is an autogenerated mock type for the Store type
type Store struct {
	mock.Mock
}

// Authenticate provides a mock function with given fields: email, password
func (_m *Store) Authenticate(email string, password string) (*model.Worker, error) {
	ret := _m.Called(email, password)

	var r0 *model.Worker
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (*model.Worker, error)); ok {
		return rf(email, password)
	}
	if rf, ok := ret.Get(0).(func(string, string) *model.Worker); ok {
		r0 = rf(email, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Worker)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(email, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateShift provides a mock function with given fields: shift
func (_m *Store) CreateShift(shift *model.Shift) error {
	ret := _m.Called(shift)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Shift) error); ok {
		r0 = rf(shift)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateShiftAssignment provides a mock function with given fields: workerId, shiftId
func (_m *Store) CreateShiftAssignment(workerId model.WorkerID, shiftId model.ShiftID) error {
	ret := _m.Called(workerId, shiftId)

	var r0 error
	if rf, ok := ret.Get(0).(func(model.WorkerID, model.ShiftID) error); ok {
		r0 = rf(workerId, shiftId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateWorker provides a mock function with given fields: worker
func (_m *Store) CreateWorker(worker *model.Worker) error {
	ret := _m.Called(worker)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Worker) error); ok {
		r0 = rf(worker)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteShiftAssignment provides a mock function with given fields: workerId, shiftId
func (_m *Store) DeleteShiftAssignment(workerId model.WorkerID, shiftId model.ShiftID) error {
	ret := _m.Called(workerId, shiftId)

	var r0 error
	if rf, ok := ret.Get(0).(func(model.WorkerID, model.ShiftID) error); ok {
		r0 = rf(workerId, shiftId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteShiftById provides a mock function with given fields: id
func (_m *Store) DeleteShiftById(id model.ShiftID) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(model.ShiftID) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteWorkerById provides a mock function with given fields: id
func (_m *Store) DeleteWorkerById(id model.WorkerID) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(model.WorkerID) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetShiftById provides a mock function with given fields: id
func (_m *Store) GetShiftById(id model.ShiftID) (*model.Shift, error) {
	ret := _m.Called(id)

	var r0 *model.Shift
	var r1 error
	if rf, ok := ret.Get(0).(func(model.ShiftID) (*model.Shift, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(model.ShiftID) *model.Shift); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Shift)
		}
	}

	if rf, ok := ret.Get(1).(func(model.ShiftID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetShifts provides a mock function with given fields: date, span, workerId
func (_m *Store) GetShifts(date *time.Time, span store.TimeSpan, workerId *model.WorkerID) ([]*model.Shift, error) {
	ret := _m.Called(date, span, workerId)

	var r0 []*model.Shift
	var r1 error
	if rf, ok := ret.Get(0).(func(*time.Time, store.TimeSpan, *model.WorkerID) ([]*model.Shift, error)); ok {
		return rf(date, span, workerId)
	}
	if rf, ok := ret.Get(0).(func(*time.Time, store.TimeSpan, *model.WorkerID) []*model.Shift); ok {
		r0 = rf(date, span, workerId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Shift)
		}
	}

	if rf, ok := ret.Get(1).(func(*time.Time, store.TimeSpan, *model.WorkerID) error); ok {
		r1 = rf(date, span, workerId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWorkerById provides a mock function with given fields: id
func (_m *Store) GetWorkerById(id model.WorkerID) (*model.Worker, error) {
	ret := _m.Called(id)

	var r0 *model.Worker
	var r1 error
	if rf, ok := ret.Get(0).(func(model.WorkerID) (*model.Worker, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(model.WorkerID) *model.Worker); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Worker)
		}
	}

	if rf, ok := ret.Get(1).(func(model.WorkerID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWorkers provides a mock function with given fields:
func (_m *Store) GetWorkers() ([]*model.Worker, error) {
	ret := _m.Called()

	var r0 []*model.Worker
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]*model.Worker, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*model.Worker); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Worker)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Migrate provides a mock function with given fields:
func (_m *Store) Migrate() {
	_m.Called()
}

// UpdateShift provides a mock function with given fields: shift
func (_m *Store) UpdateShift(shift *model.Shift) error {
	ret := _m.Called(shift)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Shift) error); ok {
		r0 = rf(shift)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateWorker provides a mock function with given fields: worker
func (_m *Store) UpdateWorker(worker *model.Worker) error {
	ret := _m.Called(worker)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Worker) error); ok {
		r0 = rf(worker)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewStore creates a new instance of Store. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStore(t mockConstructorTestingTNewStore) *Store {
	mock := &Store{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}