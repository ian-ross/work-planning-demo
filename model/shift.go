package model

import (
	"time"

	"skybluetrades.net/work-planning-demo/api"
)

type ShiftID int64

type Shift struct {
	ID        ShiftID   `db:"id"`
	StartTime time.Time `db:"start_time"`
	EndTime   time.Time `db:"end_time"`
	Capacity  int       `db:"capacity"`
}

func ShiftFromAPI(s *api.Shift) *Shift {
	var id int64
	if s.Id != nil {
		id = *s.Id
	}
	return &Shift{
		ID:        ShiftID(id),
		StartTime: s.StartTime,
		EndTime:   s.EndTime,
		Capacity:  int(s.Capacity),
	}
}

func ShiftToAPI(s *Shift) *api.Shift {
	id := int64(s.ID)
	return &api.Shift{
		Id:        &id,
		StartTime: s.StartTime,
		EndTime:   s.EndTime,
		Capacity:  int32(s.Capacity),
	}
}
