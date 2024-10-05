package models

import (
	"database/sql"
	"time"
)

// Define a User struct.  The field names and types align
// with the columns in the database "users" table.
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// Define a new UserModel struct which wraps a database connection pool.
type UserModel struct {
	DB *sql.DB
}

// The Insert method will add a new record to the "users" table.
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// This method will verify whether a user exists with the provided email
// address and password, returning the relevant user ID they do.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// This method will check if a user exists given a specific ID.
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
