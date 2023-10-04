package models

import (
	"snippetbox/internal/assert"
	"testing"
)

func TestUserModelExists(t *testing.T) {

	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	// Setup a suite of table-driven tests and expected results
	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: 1,
			want:   true,
		},
		{
			name:   "Zero ID",
			userID: 0,
			want:   false,
		},
		{
			name:   "Non-existent ID",
			userID: 2,
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the newTestDB() helper function to get a connection pool to
			// our test database.  Calling this here -- inside r.Run() -- means
			// that fresh database tables and data will be setup and torn down
			// for each subtest.
			db := newTestDB(t)

			// Create a new instance of the UserModel.
			m := UserModel{db}

			// Call the UserModel.Exists() method and check that the return
			// value and error math the expected value for the subtest.
			exists, err := m.Exists(tt.userID)
			assert.Equal(t, exists, tt.want)
			assert.NilError(t, err)
		})
	}
}
