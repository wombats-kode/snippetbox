package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls. The value of the flag
	// will be stored in the addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

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

	// The value returned from the flag.String() function is a pointer to a the flag
	// value itself, so needs to be de-referenced before using it.
	log.Printf("Starting server on %s", *addr)

	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
