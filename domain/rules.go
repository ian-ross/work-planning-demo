package domain

import "skybluetrades.net/work-planning-demo/model"

// NewShiftAssignmentOK checks the business rule:
//
//	A **worker** never has two **shifts** on the same day.
func NewShiftAssignmentOK(shifts []*model.Shift, shift *model.Shift) bool {
	// Check date of new shift against date of existing shifts.
	y, m, d := shift.StartTime.Date()
	for _, s := range shifts {
		cy, cm, cd := s.StartTime.Date()
		if y == cy && m == cm && d == cd {
			return false
		}
	}
	return true
}
