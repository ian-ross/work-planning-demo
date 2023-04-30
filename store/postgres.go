package store

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
	"golang.org/x/crypto/bcrypt"
	"skybluetrades.net/work-planning-demo/domain"
	"skybluetrades.net/work-planning-demo/model"

	// Import Postgres DB driver.
	_ "github.com/lib/pq"
)

// PGClient is a wrapper for the user database connection.
type PGClient struct {
	db *sqlx.DB
}

//go:embed postgres/*.sql
var migrations embed.FS

// NewPostgresStore creates a new user database connection.
func NewPostgresStore(dbURL string) (Store, error) {
	// Connect to database and test connection integrity.
	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	// Limit maximum connections (default is unlimited).
	db.SetMaxOpenConns(10)

	return &PGClient{db: db}, nil
}

func (pg *PGClient) Migrate() {
	// Find embedded migrations.
	m := &migrate.AssetMigrationSource{
		Asset: migrations.ReadFile,
		AssetDir: func() func(string) ([]string, error) {
			return func(path string) ([]string, error) {
				dirEntry, err := migrations.ReadDir(path)
				if err != nil {
					return nil, err
				}
				entries := make([]string, 0)
				for _, e := range dirEntry {
					entries = append(entries, e.Name())
				}

				return entries, nil
			}
		}(),
		Dir: "postgres",
	}

	// Run and log database migrations.
	_, err := migrate.Exec(pg.db.DB, "postgres", m, migrate.Up)
	if err != nil {
		log.Fatalln("Failed to migrate PostgreSQL database: ", err)
	}
}

func (pg *PGClient) Authenticate(email string, password string) (*model.Worker, error) {
	worker := &model.Worker{}
	err := pg.db.Get(worker, workerByEmail, email)
	if err == sql.ErrNoRows {
		return nil, errors.New("unknown worker email")
	}
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(worker.Password), []byte(password)); err != nil {
		return nil, err
	}
	return worker, nil
}

const workerByEmail = `
SELECT id, email, name, is_admin, password
  FROM worker
 WHERE email = $1`

func (pg *PGClient) GetWorkers() ([]*model.Worker, error) {
	results := []*model.Worker{}
	var err error
	err = pg.db.Select(&results, getWorkers)
	if err != nil {
		return nil, err
	}
	return results, nil
}

const getWorkers = `SELECT id, email, name, is_admin FROM worker`

func (pg *PGClient) GetWorkerById(id model.WorkerID) (*model.Worker, error) {
	worker := &model.Worker{}
	err := pg.db.Get(worker, workerById, id)
	if err == sql.ErrNoRows {
		return nil, errors.New("unknown worker ID")
	}
	if err != nil {
		return nil, err
	}
	return worker, nil
}

const workerById = `
SELECT id, email, name, is_admin, password
  FROM worker
 WHERE id = $1`

func (pg *PGClient) CreateWorker(worker *model.Worker) error {
	tx, err := pg.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	check := &model.Worker{}
	err = tx.Get(check, workerByEmail, worker.Email)
	if err != sql.ErrNoRows {
		return errors.New("non-unique worker email")
	}

	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(worker.Password), 0)
	stored := *worker
	stored.Password = string(bcryptPassword)

	rows, err := tx.NamedQuery(createWorker, &stored)
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		return sql.ErrNoRows
	}

	err = rows.Scan(&worker.ID)
	if err != nil {
		return err
	}

	return nil
}

const createWorker = `
INSERT INTO worker (email, name, is_admin, password)
     VALUES (:email, :name, :is_admin, :password)
RETURNING id`

func (pg *PGClient) UpdateWorker(worker *model.Worker) error {
	tx, err := pg.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(worker.Password), 0)
	stored := *worker
	stored.Password = string(bcryptPassword)

	result, err := tx.NamedExec(updateWorker, &stored)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return errors.New("failed to update worker")
	}

	return nil
}

const updateWorker = `
UPDATE worker
   SET email = :email, name = :name,
       is_admin = :is_admin, password = :password
WHERE id = :id`

func (pg *PGClient) DeleteWorkerById(id model.WorkerID) error {
	result, err := pg.db.Exec(deleteWorker, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return errors.New("failed to delete worker")
	}
	return nil
}

const deleteWorker = "DELETE FROM worker WHERE id = $1"

func (pg *PGClient) GetShifts(date *time.Time, span TimeSpan, workerId *model.WorkerID) ([]*model.Shift, error) {
	// Calculate interval start and end from date and span.
	var intStart, intEnd time.Time
	includeAll := date == nil
	if !includeAll {
		intStart, intEnd = getSpanRange(date, span)
	}

	conditions := []string{}
	if workerId != nil {
		cond := fmt.Sprintf("id = %d", *workerId)
		conditions = append(conditions, cond)
	}
	if !includeAll {
		sIntStart := intStart.Format(time.RFC3339)
		sIntEnd := intEnd.Format(time.RFC3339)
		cond := fmt.Sprintf("start_time < '%s' AND end_time > '%s'", sIntEnd, sIntStart)
		conditions = append(conditions, cond)
	}
	q := getShifts
	if len(conditions) > 0 {
		q += " WHERE " + strings.Join(conditions, " AND ")
	}

	results := []*model.Shift{}
	var err error
	err = pg.db.Select(&results, q)
	if err != nil {
		return nil, err
	}
	return results, nil
}

const getShifts = `SELECT id, start_time, end_time, capacity FROM shift`

func (pg *PGClient) GetShiftById(id model.ShiftID) (*model.Shift, error) {
	shift := &model.Shift{}
	err := pg.db.Get(shift, shiftById, id)
	if err == sql.ErrNoRows {
		return nil, errors.New("unknown shift ID")
	}
	if err != nil {
		return nil, err
	}
	return shift, nil
}

const shiftById = `
SELECT id, start_time, end_time, capacity
  FROM shift
 WHERE id = $1`

func (pg *PGClient) CreateShift(shift *model.Shift) error {
	tx, err := pg.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	rows, err := tx.NamedQuery(createWorker, shift)
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		return sql.ErrNoRows
	}

	err = rows.Scan(&shift.ID)
	if err != nil {
		return err
	}

	return nil
}

const createShift = `
INSERT INTO shift (start_time, end_time, capacity)
     VALUES (:start_time, :end_time, :capacity)
RETURNING id`

func (pg *PGClient) UpdateShift(shift *model.Shift) error {
	tx, err := pg.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	result, err := tx.NamedExec(updateShift, shift)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return errors.New("failed to update shift")
	}

	return nil
}

const updateShift = `
UPDATE worker
   SET start_time = :start_time, end_time = :end_time,
       capacity = :capacity
WHERE id = :id`

func (pg *PGClient) DeleteShiftById(id model.ShiftID) error {
	result, err := pg.db.Exec(deleteShift, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return errors.New("failed to delete shift")
	}
	return nil
}

const deleteShift = "DELETE FROM shift WHERE id = $1"

func (pg *PGClient) CreateShiftAssignment(workerId model.WorkerID, shiftId model.ShiftID) error {
	tx, err := pg.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	worker := &model.Worker{}
	err = tx.Get(worker, workerById, workerId)
	if err == sql.ErrNoRows {
		return errors.New("unknown worker ID")
	}
	if err != nil {
		return err
	}

	shift := &model.Shift{}
	err = tx.Get(shift, shiftById, shiftId)
	if err == sql.ErrNoRows {
		return errors.New("unknown shift ID")
	}
	if err != nil {
		return err
	}

	shifts, err := pg.GetShifts(nil, WeekSpan, &workerId)
	if err != nil {
		return ErrRetrievingWorkerShifts
	}

	if !domain.NewShiftAssignmentOK(shifts, shift) {
		return ErrTwoShiftsSameDay
	}

	assignment := &model.ShiftAssignment{Worker: workerId, Shift: shiftId}
	rows, err := tx.NamedQuery(createShiftAssignment, assignment)
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		return sql.ErrNoRows
	}

	return nil
}

const createShiftAssignment = `
INSERT INTO shift_assignment (worker_id, shift_id)
     VALUES (:worker_id, :shift_id)
RETURNING id`

func (pg *PGClient) DeleteShiftAssignment(workerId model.WorkerID, shiftId model.ShiftID) error {
	result, err := pg.db.Exec(deleteShiftAssignment, workerId, shiftId)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return errors.New("failed to delete shift assignment")
	}
	return nil
}

const deleteShiftAssignment = `
DELETE FROM shift_assignment
 WHERE worker_id = $1 AND shift_id = $2`
