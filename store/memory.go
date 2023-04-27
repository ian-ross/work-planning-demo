package store

import (
	"errors"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
	"skybluetrades.net/work-planning-demo/model"
)

type MemoryStore struct {
	sync.RWMutex
	lastWorkerID     model.WorkerID
	lastShiftID      model.ShiftID
	lastAssignmentID model.ShiftAssignmentID
	workers          map[model.WorkerID]*model.Worker
	workersByEmail   map[string]*model.Worker
	shifts           map[model.ShiftID]*model.Shift
	assignments      map[model.ShiftAssignmentID]*model.ShiftAssignment
}

func NewMemoryStore() (Store, error) {
	return &MemoryStore{
		lastWorkerID:     0,
		lastShiftID:      0,
		lastAssignmentID: 0,
		workers:          make(map[model.WorkerID]*model.Worker),
		workersByEmail:   make(map[string]*model.Worker),
		shifts:           make(map[model.ShiftID]*model.Shift),
		assignments:      make(map[model.ShiftAssignmentID]*model.ShiftAssignment),
	}, nil
}

func (s *MemoryStore) Authenticate(email string, password string) (*model.Worker, error) {
	s.RLock()
	defer s.RUnlock()

	worker, exists := s.workersByEmail[email]
	if !exists {
		return nil, errors.New("unknown user email")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(worker.Password), []byte(password)); err != nil {
		return nil, err
	}

	retWorker := *s.workers[worker.ID]
	return &retWorker, nil
}

func (s *MemoryStore) GetWorkers() ([]*model.Worker, error) {
	s.RLock()
	s.RUnlock()

	workers := make([]*model.Worker, len(s.workers))
	i := 0
	for _, w := range s.workers {
		rworker := *w
		workers[i] = &rworker
		i++
	}

	return workers, nil
}

func (s *MemoryStore) GetWorkerById(id model.WorkerID) (*model.Worker, error) {
	s.RLock()
	defer s.RUnlock()

	worker, exists := s.workers[id]
	if !exists {
		return nil, errors.New("unknown worker ID")
	}

	rworker := *worker
	return &rworker, nil
}

func (s *MemoryStore) CreateWorker(worker *model.Worker) (*model.Worker, error) {
	s.RLock()
	defer s.RUnlock()

	stored := *worker
	s.lastWorkerID++
	stored.ID = s.lastWorkerID
	s.workers[stored.ID] = &stored
	s.workersByEmail[stored.Email] = &stored

	worker.ID = stored.ID
	return worker, nil
}

func (s *MemoryStore) UpdateWorker(worker *model.Worker) (*model.Worker, error) {
	s.RLock()
	defer s.RUnlock()

	existing, exists := s.workers[worker.ID]
	if !exists {
		return nil, errors.New("unknown worker ID")
	}

	delete(s.workers, worker.ID)
	delete(s.workersByEmail, existing.Email)

	stored := *worker
	s.workers[stored.ID] = &stored
	s.workersByEmail[stored.Email] = &stored

	return worker, nil
}

func (s *MemoryStore) DeleteWorker(worker *model.Worker) error {
	s.RLock()
	defer s.RUnlock()

	existing, exists := s.workers[worker.ID]
	if !exists {
		return errors.New("unknown worker ID")
	}

	delete(s.workers, worker.ID)
	delete(s.workersByEmail, existing.Email)

	return nil
}

func (s *MemoryStore) GetShifts(date *time.Time, span TimeSpan) ([]*model.Shift, error) {
	s.RLock()
	s.RUnlock()

	// Calculate interval start and end from date and span.
	y, m, d := date.Date()
	intStart := time.Date(y, m, d, 0, 0, 0, 0, nil)
	var intEnd time.Time
	if span == DaySpan {
		intEnd = intStart.AddDate(0, 0, 1)
	} else {
		wd := intStart.Weekday()
		delta := (int(wd) - 1) % 7
		intStart = intStart.AddDate(0, 0, -delta)
		intEnd = intStart.AddDate(0, 0, 7)
	}

	// Include only shifts in interval.
	shifts := make([]*model.Shift, len(s.shifts))
	i := 0
	for _, s := range s.shifts {
		if s.StartTime.Before(intEnd) && s.EndTime.After(intStart) {
			rshift := *s
			shifts[i] = &rshift
			i++
		}
	}

	return shifts, nil
}

func (s *MemoryStore) GetShiftById(id model.ShiftID) (*model.Shift, error) {
	s.RLock()
	defer s.RUnlock()

	shift, exists := s.shifts[id]
	if !exists {
		return nil, errors.New("unknown shift ID")
	}

	rshift := *shift
	return &rshift, nil
}

func (s *MemoryStore) CreateShift(shift *model.Shift) (*model.Shift, error) {
	s.RLock()
	defer s.RUnlock()

	stored := *shift
	s.lastShiftID++
	stored.ID = s.lastShiftID
	s.shifts[stored.ID] = &stored

	shift.ID = stored.ID
	return shift, nil
}

func (s *MemoryStore) UpdateShift(shift *model.Shift) (*model.Shift, error) {
	s.RLock()
	defer s.RUnlock()

	_, exists := s.shifts[shift.ID]
	if !exists {
		return nil, errors.New("unknown shift ID")
	}

	delete(s.shifts, shift.ID)

	stored := *shift
	s.shifts[stored.ID] = &stored

	return shift, nil
}

func (s *MemoryStore) DeleteShift(shift *model.Shift) error {
	s.RLock()
	defer s.RUnlock()

	_, exists := s.shifts[shift.ID]
	if !exists {
		return errors.New("unknown shift ID")
	}

	delete(s.shifts, shift.ID)

	return nil
}

func (s *MemoryStore) CreateShiftAssignment(
	workerId model.WorkerID, shiftId model.ShiftID) (*model.ShiftAssignment, error) {
	return nil, errors.New("NYI: MemoryStore.CreateShiftAssignment")
}

func (s *MemoryStore) DeleteShiftAssignment(assignment *model.ShiftAssignment) error {
	return errors.New("NYI: MemoryStore.DeleteShiftAssignment")
}
