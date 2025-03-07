package models

import (
	"database/sql"
	"errors"
	"fmt"
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

	_, err = m.DB.Exec(stmt, name, email, password)
	if err != nil {
		var pqErr *pq.Error
		errors.As(err, &pqErr)
		// unique_violation 23505
		fmt.Printf("Class, name, code: %#v, %#v, %#v\n", pqErr.Code.Class(), pqErr.Code.Name(), pqErr.Code)
		if pqErr.Code.Name() == "unique_violation" {
			fmt.Printf("INSERT ERROR UNIQUE_VIOLATION", pqErr)
			return ErrDuplicateEmail
		}

		return err
	}

	fmt.Printf("ALL GOOD: %#v\n", hashedPassword)
	return nil
}

func (*UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (*UserModel) Exists(id int) (bool, error) {
	return false, nil
}
