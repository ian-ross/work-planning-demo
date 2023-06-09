// `MemoryStore` is a simple in-memory database for workers and
// shifts. It implements the `Store` interface, and includes some test
// data setup via its `Migrate` method (called from the main program
// during application startup).

package store

import (
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slices"
	"skybluetrades.net/work-planning-demo/domain"
	"skybluetrades.net/work-planning-demo/model"
)

type MemoryStore struct {
	sync.RWMutex
	lastWorkerID   model.WorkerID
	lastShiftID    model.ShiftID
	workers        map[model.WorkerID]*model.Worker
	workersByEmail map[string]*model.Worker
	shifts         map[model.ShiftID]*model.Shift
	assignments    []model.ShiftAssignment
}

func NewMemoryStore() (Store, error) {
	return &MemoryStore{
		lastWorkerID:   0,
		lastShiftID:    0,
		workers:        make(map[model.WorkerID]*model.Worker),
		workersByEmail: make(map[string]*model.Worker),
		shifts:         make(map[model.ShiftID]*model.Shift),
		assignments:    []model.ShiftAssignment{},
	}, nil
}

func (s *MemoryStore) Migrate() {
	if len(s.workers) == 0 && len(s.shifts) == 0 && len(s.assignments) == 0 {
		addTestData(s)
	}
}

func (s *MemoryStore) Authenticate(email string, password string) (*model.Worker, error) {
	s.RLock()
	defer s.RUnlock()

	worker, exists := s.workersByEmail[email]
	if !exists {
		return nil, ErrUnknownWorkerEmail
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
		return nil, ErrWorkerNotFound
	}

	rworker := *worker
	return &rworker, nil
}

func (s *MemoryStore) CreateWorker(worker *model.Worker) error {
	s.RLock()
	defer s.RUnlock()

	stored := *worker
	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(worker.Password), 0)
	stored.Password = string(bcryptPassword)
	s.lastWorkerID++
	stored.ID = s.lastWorkerID
	s.workers[stored.ID] = &stored
	s.workersByEmail[stored.Email] = &stored

	worker.ID = stored.ID
	return nil
}

func (s *MemoryStore) UpdateWorker(worker *model.Worker) error {
	s.RLock()
	defer s.RUnlock()

	existing, exists := s.workers[worker.ID]
	if !exists {
		return ErrWorkerNotFound
	}

	delete(s.workers, worker.ID)
	delete(s.workersByEmail, existing.Email)

	stored := *worker
	s.workers[stored.ID] = &stored
	s.workersByEmail[stored.Email] = &stored

	return nil
}

func (s *MemoryStore) DeleteWorkerById(id model.WorkerID) error {
	s.RLock()
	defer s.RUnlock()

	existing, exists := s.workers[id]
	if !exists {
		return ErrWorkerNotFound
	}

	delete(s.workers, id)
	delete(s.workersByEmail, existing.Email)

	return nil
}

func (s *MemoryStore) GetShifts(
	date *time.Time, span TimeSpan, workerId *model.WorkerID) ([]*model.Shift, error) {
	s.RLock()
	s.RUnlock()

	// Calculate interval start and end from date and span.
	var intStart, intEnd time.Time
	includeAll := date == nil
	if !includeAll {
		intStart, intEnd = getSpanRange(date, span)
	}

	// If we're extracting shifts for a given worker, collect the
	// assigned shift IDs for the worker for filtering here.
	var assigned []model.ShiftID
	if workerId != nil {
		for _, a := range s.assignments {
			if a.Worker == *workerId {
				assigned = append(assigned, a.Shift)
			}
		}
	}

	// Include only shifts in interval.
	shifts := []*model.Shift{}
	i := 0
	for _, s := range s.shifts {
		include := includeAll || s.StartTime.Before(intEnd) && s.EndTime.After(intStart)
		if workerId != nil && include {
			include = slices.Contains(assigned, s.ID)
		}
		if include {
			rshift := *s
			shifts = append(shifts, &rshift)
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
		return nil, ErrShiftNotFound
	}

	rshift := *shift
	return &rshift, nil
}

func (s *MemoryStore) CreateShift(shift *model.Shift) error {
	s.RLock()
	defer s.RUnlock()

	stored := *shift
	s.lastShiftID++
	stored.ID = s.lastShiftID
	s.shifts[stored.ID] = &stored

	shift.ID = stored.ID
	return nil
}

func (s *MemoryStore) UpdateShift(shift *model.Shift) error {
	s.RLock()
	defer s.RUnlock()

	_, exists := s.shifts[shift.ID]
	if !exists {
		return ErrShiftNotFound
	}

	delete(s.shifts, shift.ID)

	stored := *shift
	s.shifts[stored.ID] = &stored

	return nil
}

func (s *MemoryStore) DeleteShiftById(id model.ShiftID) error {
	s.RLock()
	defer s.RUnlock()

	_, exists := s.shifts[id]
	if !exists {
		return ErrShiftNotFound
	}

	delete(s.shifts, id)

	return nil
}

func (s *MemoryStore) CreateShiftAssignment(
	workerId model.WorkerID, shiftId model.ShiftID) error {
	s.RLock()
	defer s.RUnlock()

	shift, exists := s.shifts[shiftId]
	if !exists {
		return ErrShiftNotFound
	}

	existing := 0
	for _, a := range s.assignments {
		if a.Shift == shiftId {
			existing++
		}
	}
	if existing >= shift.Capacity {
		return ErrShiftAtCapacity
	}

	shifts, err := s.GetShifts(nil, WeekSpan, &workerId)
	if err != nil {
		return ErrRetrievingWorkerShifts
	}

	if !domain.NewShiftAssignmentOK(shifts, shift) {
		return ErrTwoShiftsSameDay
	}

	s.assignments = append(s.assignments, model.ShiftAssignment{Worker: workerId, Shift: shiftId})
	return nil
}

func (s *MemoryStore) DeleteShiftAssignment(
	workerId model.WorkerID, shiftId model.ShiftID) error {
	s.RLock()
	defer s.RUnlock()

	pos := slices.Index(s.assignments, model.ShiftAssignment{Worker: workerId, Shift: shiftId})
	if pos == -1 {
		return ErrShiftAssignmentNotFound
	}

	s.assignments = slices.Delete(s.assignments, pos, pos)
	return nil
}
