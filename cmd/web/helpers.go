package main

import (
	"net/http"
	"runtime/debug"
)

// The ServerError helper writes a log entry at Error level (including the request
// method and URI attributes), then sends a generic 500 internal Server Error
// response to the user.
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack()) // Convert stack trace from a []byte to String
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError help sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// request" when there's a problem with the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
