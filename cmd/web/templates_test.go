package main

import (
	"snippetbox/internal/assert"
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {

	// create a slice of anonymous structs to containing the test case name,
	// input to our humanDate() function (the tm field), and expected output
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "17 Mar 2022 at 10:15",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2022, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Mar 2022 at 09:15",
		},
	}

	// Loop over the test cases.
	for _, tt := range tests {
		// Use the t.Run() function to run a subtest for each test case.  The
		// first parameter to this is the name, and the second parameter is an
		// anonymous function containing the actual test for each case.
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			// Use the new assert.Equal() function to simplify the test
			assert.Equal(t, hd, tt.want)
		})

	}
}
