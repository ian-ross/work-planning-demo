package server

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"skybluetrades.net/work-planning-demo/model"
	"skybluetrades.net/work-planning-demo/store"
)

func createWorker(s store.Store,
	email string, name string, password string, isAdmin bool) *model.Worker {
	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 0)
	worker, err := s.CreateWorker(&model.Worker{
		Email:    email,
		Name:     name,
		IsAdmin:  isAdmin,
		Password: string(bcryptPassword),
	})
	if err != nil {
		log.Fatalln("Failed creating test data in createWorker: ", err)
	}
	return worker
}

func createShift(s store.Store,
	day string, hourStart int, capacity int) *model.Shift {
	dayDate, _ := time.Parse("2006-01-02", day)
	start := dayDate.Add(time.Duration(hourStart) * time.Hour)
	end := dayDate.Add(time.Duration(hourStart+8) * time.Hour)
	shift, err := s.CreateShift(&model.Shift{
		StartTime: start,
		EndTime:   end,
		Capacity:  capacity,
	})
	if err != nil {
		log.Fatalln("Failed creating test data in createShift: ", err)
	}
	return shift
}

func addTestData(s store.Store) {
	fmt.Println()
	fmt.Println("+---------------------+")
	fmt.Println("| ADDING TEST DATA... |")
	fmt.Println("+---------------------+")
	fmt.Println()

	workers := []*model.Worker{}
	workers = append(workers, createWorker(s, "test1@example.com", "Tina Tester", "password1", true))
	workers = append(workers, createWorker(s, "test2@example.com", "Tom Testerman", "password2", false))
	workers = append(workers, createWorker(s, "test3@example.com", "Tammy Testino", "password3", false))
	workers = append(workers, createWorker(s, "test4@example.com", "Todd Testa", "password4", false))

	shifts := []*model.Shift{}
	for d := 1; d <= 7; d++ {
		dt := fmt.Sprintf("2023-05-0%d", d)
		shifts = append(shifts, createShift(s, dt, 0, 1))
		shifts = append(shifts, createShift(s, dt, 8, 3))
		shifts = append(shifts, createShift(s, dt, 16, 2))
	}
}
