package model

type WorkerID int64

type Worker struct {
	ID       WorkerID `db:"id"`
	Email    string   `db:"email"`
	Name     string   `db:"name"`
	IsAdmin  bool     `db:"is_admin"`
	Password string   `db:"password"`
}
