package store

import (
	"errors"
)

// PGClient is a wrapper for the user database connection.
type PGClient struct {
	// DB *sqlx.DB
}

// NewPostgresStore creates a new user database connection.
func NewPostgresStore(dbURL string) (Store, error) {
	return nil, errors.New("TBD")
}
