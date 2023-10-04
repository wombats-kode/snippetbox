package models

import (
	"database/sql"
	"os"
	"testing"
)

func newTestDB(t *testing.T) *sql.DB {
	// Establish a sql.DB connnection pool for out test database.  Because our
	// setup and teardown scripts contain multiple sql statements, we need to
	// use the 'multistatments=true' parameter in our DSN.  This instructs
	// our MSQL database driver to support executing multiple sql statements
	// in one db.Exec() call.
	db, err := sql.Open("mysql", "test_web:pass@/test_snippetbox?parseTime=true&multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}

	// Read the setup SQL script from file and execute the statements.
	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	// Use the t.Cleanup() to register a function which will automatically be
	// called by Go when the current test (or sub-test) which call the newTestDb()
	// has finised.  In this function we read and run the Teardown script and close]
	// the db connection.
	t.Cleanup(func() {
		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	})
	return db
}
