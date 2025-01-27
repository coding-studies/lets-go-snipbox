package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

// Insert inserts a new snippet into the snippets table.
//
// The returned int is the ID of the newly inserted snippet (unless there was an
// error, in which case 0 is returned).
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	exp := fmt.Sprintf("%d DAYS", expires)

	stmt := `
		INSERT INTO snippets (
			  title
			, content
			, created
			, expires
) VALUES (
			$1
		, $2
		, CURRENT_TIMESTAMP
		, CURRENT_TIMESTAMP + $3::INTERVAL
) RETURNING id;
`

	res := m.DB.QueryRow(stmt, title, content, exp)

	var lastInsertedID int
	err := res.Scan(&lastInsertedID)
	if err != nil {
		return 0, err
	}

	return lastInsertedID, nil
}

func (m *SnippetModel) Get(id int) (Snippet, error) {
	return Snippet{}, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
