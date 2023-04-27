package model

import "time"

type ShiftID int64

type Shift struct {
	ID        ShiftID   `db:"id"`
	StartTime time.Time `db:"start_time"`
	EndTime   time.Time `db:"end_time"`
	Capacity  int       `db:"capacity"`
}
