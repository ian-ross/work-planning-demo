package model

type ShiftAssignmentID int64

type ShiftAssignment struct {
	ID     ShiftAssignmentID `db:"id"`
	Worker WorkerID          `db:"worker_id"`
	Shift  ShiftID           `db:"shift_id"`
}
