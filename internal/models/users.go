package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `
		INSERT INTO users (
				name
			, email
			, hashed_password
			, created
		) VALUES(
				$1
			, $2
			, $3
			, CURRENT_TIMESTAMP
		);
	`

	_, err = m.DB.Exec(stmt, name, email, hashedPassword)
	if err != nil {
		var pqErr *pq.Error
		errors.As(err, &pqErr)
		// unique_violation 23505
		if pqErr.Code.Name() == "unique_violation" {
			return ErrDuplicateEmail
		}

		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = $1"

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return 0, nil
}

func (*UserModel) Exists(id int) (bool, error) {
	return false, nil
}
