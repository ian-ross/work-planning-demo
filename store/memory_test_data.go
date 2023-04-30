package store

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"skybluetrades.net/work-planning-demo/model"
)

func createTestWorker(s Store,
	email string, name string, password string, isAdmin bool) *model.Worker {
	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 0)
	worker := &model.Worker{
		Email:    email,
		Name:     name,
		IsAdmin:  isAdmin,
		Password: string(bcryptPassword),
	}
	err := s.CreateWorker(worker)
	if err != nil {
		log.Fatalln("Failed creating test data in createWorker: ", err)
	}
	return worker
}

func createTestShift(s Store,
	day string, hourStart int, capacity int) *model.Shift {
	dayDate, _ := time.Parse("2006-01-02", day)
	start := dayDate.Add(time.Duration(hourStart) * time.Hour)
	end := dayDate.Add(time.Duration(hourStart+8) * time.Hour)
	shift := &model.Shift{
		StartTime: start,
		EndTime:   end,
		Capacity:  capacity,
	}
	err := s.CreateShift(shift)
	if err != nil {
		log.Fatalln("Failed creating test data in createShift: ", err)
	}
	return shift
}

func addTestData(s Store) {
	fmt.Println()
	fmt.Println("+---------------------+")
	fmt.Println("| ADDING TEST DATA... |")
	fmt.Println("+---------------------+")
	fmt.Println()

	workers := []*model.Worker{}
	workers = append(workers, createTestWorker(s, "test1@example.com", "Tina Tester", "password1", true))
	workers = append(workers, createTestWorker(s, "test2@example.com", "Tom Testerman", "password2", false))
	workers = append(workers, createTestWorker(s, "test3@example.com", "Tammy Testino", "password3", false))
	workers = append(workers, createTestWorker(s, "test4@example.com", "Todd Testa", "password4", false))

	shifts := []*model.Shift{}
	for d := 1; d <= 31; d++ {
		dt := fmt.Sprintf("2023-05-%02d", d)
		ndaytime := 3
		if d%7 == 6 || d%7 == 0 {
			// Weekends
			ndaytime = 2
		}
		shifts = append(shifts, createTestShift(s, dt, 0, 1))
		shifts = append(shifts, createTestShift(s, dt, 8, ndaytime))
		shifts = append(shifts, createTestShift(s, dt, 16, 2))
	}

	for i := 0; i < 7; i++ {
		s.CreateShiftAssignment(workers[1].ID, shifts[i*3+1].ID)
	}
}
