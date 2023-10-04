package models

import (
	"database/sql"
	"errors"
	"time"
)

type SnippetModelInterface interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (*Snippet, error)
	Latest() ([]*Snippet, error)
}

// Define a snippet type to hold the data for an individual snippet. Notice how
// the fields of the struct correspond to the fields in our MySQL snippets
// table.
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a SnipperModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// Function to insert a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Use the Exec() method on the connection pool to execute the statement.
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Use the LastInsertId() method on the result to get the ID of our newley
	// created record in the snippets table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// The ID retunred has the type int64, so we convert it to an int type
	// before returning
	return int(id), nil
}

// Function that will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires
	FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// USe the QueryRow method on the connection pool to execute this statement.
	// This returns a pointer to a sql.Row object which holds the result from
	// the database
	row := m.DB.QueryRow(stmt, id)

	// Create a pointer to a new zeroed snippet struct.
	s := &Snippet{}

	// Use row.Scan() to copy the values from each field to the corresponding
	//field in the Snippet struct.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// if the query returns no rows, then row.Scan() will return a
		// sql.ErrNoRows error.  We use errors.Is() function to check for
		// that specific error and return our own ErrNoRecord.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	// If everthing went OK the return the Snippet object
	return s, nil
}

// Function to return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {

	stmt := `SELECT id, title, content, created, expires
	FROM snippets
	WHERE expires > UTC_TIMESTAMP()
	ORDER BY id DESC
	LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// We defer rows.Close to ensure the sql.Rows resultset is
	// always properly closed before the Latest() method returns.
	defer rows.Close()

	// Initialise an empty slice to hold the Snippet structs
	snippets := []*Snippet{}

	// Use rows.Next() to iterate through the rows in the resultset.
	for rows.Next() {
		// Create a pointer to a new zeroed Snippet struct
		s := &Snippet{}

		// Use rows.Scan() to copy the values from each field in the row to the
		// new Snippet object.  The arguments must be pointers to the place where
		// the data is to be copied and the number of arguments must be exactly the
		// same as the number fo columns returned by the statement.
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	// When the rows.Next() loop has finished we call rows.Err() to retrieve any error
	// that was encoutered duriing the interation.  Its important to call this - do not
	// assume that succeddful interation was completed.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// If everthing went OK then return the Snippets slice

	return snippets, nil
}
