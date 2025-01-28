package models

import (
	"database/sql"
	"errors"
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
	stmt := `
		SELECT
				id
			, title
			, content
			, created
			, expires
		FROM snippets
		WHERE
			expires > CURRENT_TIMESTAMP
		AND
			id = $1;
`

	var snippet Snippet
	err := m.DB.QueryRow(stmt, id).Scan(
		&snippet.ID,
		&snippet.Title,
		&snippet.Content,
		&snippet.Created,
		&snippet.Expires,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	return snippet, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error) {
	stmt := `
		SELECT
				id
			, title
			, content
			, created
			, expires
	FROM snippets
	WHERE
		expires > CURRENT_TIMESTAMP
	ORDER BY id DESC
	LIMIT 10;
`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snippets []Snippet

	for rows.Next() {
		var snippet Snippet

		err := rows.Scan(
			&snippet.ID,
			&snippet.Title,
			&snippet.Content,
			&snippet.Created,
			&snippet.Expires,
		)

		if err != nil {
			return nil, err
		}

		snippets = append(snippets, snippet)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
