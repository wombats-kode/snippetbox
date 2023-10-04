package models

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching record found")

	// we'll use this if a user tries to login with an incorrect email
	// address or password
	ErrInvalidCredentials = errors.New("models: invalid credentials")

	// we use this if a user tries to signup with an email that is
	// already in use.
	ErrDuplicateEmail = errors.New("models: duplicate email")
)
