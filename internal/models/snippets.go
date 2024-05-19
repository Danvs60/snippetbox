package models

import (
	"database/sql"
	"errors"
	"time"
)

// Define a Snippet type mapping to the database
// fields for snippets
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Wrapper for a sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// Database commands
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// SQL statement to insert snippets.
	// use backticks to define the string in multiple lines
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Use exec method from embedded connection pool for 
	// non-query (not SELECT) statements
	result, err := m.DB.Exec(stmt, title, content, expires)
	// NOTE: can also ignore results 
	// _, err := ...
	if err != nil {
		return 0, err
	}

	// get id of last inserted snippet
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// could also use RowsAffected() to get number of rows affected

	// WARNING: not all drivers support such result functions
	// for example PostgreSQL does not support LastInsertId()

	// return the id, cast from int64 to int
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() and id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &Snippet{}

	// scan only accepts pointers (mem. addresses) as input fields
	// also the number of pointer parameters given to Scan
	// will need to exactly match the number of columns given by the statement
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	// NOTE: this will take the query raw input and map it to Go standard types
	// CHAR, VARCHAR and TEXT map to string
	// BOOLEAN maps to bool
	// INT maps to int
	// BIGINT maps to int64
	// DECIMAL and NUMERIC map to float
	// TIMESTAMP map to time.Time
	// and thanks to parseTime=true TIME, DATE also map to time.Time (otherwise they would map to []byte

	if err != nil {
		// check if error is 0 rows returned
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// should defer *after* checking for an error coming from Query()
	// otherwise, if Query returns an error, it will panic as you'd be
	// trying to close a nil resultset
	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	// check for errors during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// all good
	return snippets, nil
}
