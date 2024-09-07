package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Create a fileserver which serves files out of the "./ui/static" directory
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileserver := http.FileServer(http.Dir("./ui/static"))

	// USe the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/".  For matching paths, we strip the
	// "/static/" prefeix before the request reaches the file server.
	mux.Handle("GET /static/", http.StripPrefix("/static", fileserver))

	// Register the other application routes as normal...
	mux.HandleFunc("GET /{$}", home)                      // Restrict this route to exact matches only
	mux.HandleFunc("GET /snippet/view/{id}", snippetView) // add the {id} wildcard
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	log.Print("Starting server on :4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
