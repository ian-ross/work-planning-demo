package model

import (
	"strings"

	"skybluetrades.net/work-planning-demo/api"
)

type WorkerID int64

type Worker struct {
	ID       WorkerID `db:"id"`
	Email    string   `db:"email"`
	Name     string   `db:"name"`
	IsAdmin  bool     `db:"is_admin"`
	Password string   `db:"password"`
}

func WorkerFromAPI(w *api.Worker) *Worker {
	var id int64
	if w.Id != nil {
		id = *w.Id
	}
	return &Worker{
		ID:      WorkerID(id),
		Email:   strings.TrimSpace(w.Email),
		Name:    strings.TrimSpace(w.Name),
		IsAdmin: w.IsAdmin,
	}
}

func WorkerToAPI(worker *Worker) *api.Worker {
	id := int64(worker.ID)
	return &api.Worker{
		Id:      &id,
		Email:   worker.Email,
		Name:    worker.Name,
		IsAdmin: worker.IsAdmin,
	}
}
