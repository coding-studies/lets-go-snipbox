package models

import (
	"database/sql"
	"time"
)

// User represents a user, which maps fields from the database users table.
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// UserModel wraps our DB connection pool
type UserModel struct {
	DB *sql.DB
}

func (*UserModel) Insert(name, email, password string) error {
	return nil
}

func (*UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (*UserModel) Exists(id int) (bool, error) {
	return false, nil
}