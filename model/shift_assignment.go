package model

type ShiftAssignment struct {
	Worker WorkerID `db:"worker_id"`
	Shift  ShiftID  `db:"shift_id"`
}
