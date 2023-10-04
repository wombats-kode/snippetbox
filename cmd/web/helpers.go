package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
)

// The serverError helper writes an error message and stack trace to the errorlog
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
// to the user.  We'll use this to send responses like 400 "Bad Request" when there's a
// problem with the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// For consistency, we'll implement a notFound helper.  This is simplt a
// convenience wrapper around clientError which sends a 404 Not Found response
// to the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	// Retrieve the appropriate template set from from the cache based on the page
	// name like 'home.tmpl'.  If not entry exists in the cache with the provided name,
	// then create a new error and call the serverError() helper method that we created
	// earlier and return.
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// initialise a new buffer
	buf := new(bytes.Buffer)

	// Write the template to the buffer instead of straighe to the
	// http.responseWriter. If there is an error, call our serverError() helper
	// and then return
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// if the template is written to the buffer without any errors, we are safe
	// to go ahead and write the HTTP status code to the http.responseWriter.
	w.WriteHeader(status)

	// Write the contents of the buffer directly to the http.responseWriter.  Note
	// this is another instance where we pass our http.responseWriter to a function
	// that take an io.writer.
	buf.WriteTo(w)
}

// Create a newTemplateData() helper, which returns a pointer to a templateData
// struct initialised with the current year.
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r), // Add the CRSF token.
	}
}

// Create a new decodePostForm() helper method. The second parameter here, dst,
// is the target destination that we want to decode the form data into.
func (app *application) decodePostForm(r *http.Request, dst any) error {
	// Call ParseForm() on the requestm in the same way that we did in our
	// createSnippetPost handler
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// Call Decode() on our decoder instance, passing the target destination as
	// the first parameter.
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// If we try to use an invalid target destinationm the Decode() method
		// will return an error with the type *form.InvalidDecoderError. We use
		// errors.As() to check for this and raise a panic rather then returning
		// an error.
		var InvalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &InvalidDecoderError) {
			panic(err)
		}
		// for all other errorsm we return them as normal
		return err
	}
	return nil
}

// func returns true if the current request is from an authenticated user,
// otherwise it returns false.
func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}
