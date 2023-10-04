package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// Middleware that adds OWASP recommended headers to all ServeMux response headers
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Note this is split across multiple lines for readability.  This is not
		// required if you are after more succinct code.
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function whcih will always be run in the event of
		// a panic as Go unwinds the stack.
		defer func() {
			// Use the builtin recover function to check if there has been a panic
			// or not. If there has...
			if err := recover(); err != nil {
				// Set a connection: close header on the response
				w.Header().Set("Connection", "close")
				// Call the app.serverError helper to return a 500 Internal Server response.
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If a user is not authenticated, redirect them to the login page and
		// return from the middleware chain so that no subsequent handlers in
		// the chain are executed.
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		// Otherwise set the 'Cache-Control: no-store" header so the pages
		// require authentication are not stored in the users browser cache,
		// or other intermediary cache.
		w.Header().Add("Cache-Control", "no-store")

		// and call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

// Create a NoSurf middleware function which uses a customised CSRF cookie with
// the Secure, Path and HttpOnly attributes set.
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// retrieve the authenticatedUserID value from the session using the GetInt()
		// method.  This will return the zero value of an int if no "authenticatedUserID"
		// value is in the session -- in which case we call the next handler in the chain
		// as normal and return.
		id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		// Otherwise we check to see if a user with that ID exists in our Database
		exists, err := app.users.Exists(id)
		if err != nil {
			app.serverError(w, err)
			return
		}

		// If a matching user is found, we know that the request is coming from
		// an authenticated user who exists in the database.  We create a new copy
		// of the reqeust and assign it to r.
		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}
		// Call the next server in the chain
		next.ServeHTTP(w, r)
	})
}
