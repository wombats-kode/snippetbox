package main

import "net/http"

// The routes() method returns a servemux containing our application routes.
// Update the signature for the routes() method so that is returns a http.Handler
// instead of *http.ServeMux to support our middleware patterns.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileserver := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileserver))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	// Pass the serveMux as the 'next' parameter to the commonHeaders middleware
	// Because commonHeader is jut a function and returns a http.Hander we dont
	// need to anything else.
	return app.recoverPanic(app.logRequest(commonHeaders(mux)))
}
