package main

import (
	"net/http"
	"snippetbox/ui"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// The routes() method has been changed to return an http.Handler instead of an *http.ServeMux
// so that we can incorporate our idiomatic middlesware design to the application.
func (app *application) routes() http.Handler {
	router := httprouter.New()

	// Create a handler function that wraps our notFound() helper, and then assigns it
	// as the custom handler for 404 Not Found response.  We can also set a custom handler
	// for the 405 Method Not Allowed responses by setting router.MethoNotAllowed in the
	// same way too.
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// Take the ui.Files embedded filesystem and convert it to a http.FS type so
	// that it satisfies the http.filesystem interface.  We then pass that to the
	// http.FileServer() function to create the file server handler.
	fileserver := http.FileServer(http.FS(ui.Files))

	// Using embedded files means that strip prefix is no longer needed for our static files
	// from this: router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileserver))
	router.Handler(http.MethodGet, "/static/*filepath", fileserver)

	// Add a net Get /ping route for our testing
	router.HandlerFunc(http.MethodGet, "/ping", ping)

	// Unprotected application routes using the 'dynamic' middleware chain
	// Use the nosurf middleware on all our 'dynamic' routes.
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	// Protected(authentication-only) application routes using a new 'protected' middleware chain
	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	// Create a middleware chain containing our 'standard' middleware which will be used
	// for every request our application will receive.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// return the 'standard' middleware chain followed by the servemux.
	return standard.Then(router)

}
