package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"snippetbox.example.com/internal/models"
	"snippetbox.example.com/ui"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
type templateData struct {
	CurrentYear     int
	Snippet         models.Snippet
	Snippets        []models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

// Function that returns a nicely formatted string representation of
// a time.Time object.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// Initialise a template.FuncMap object and store it in a global variable.
// This is essentially a string-keyed map which acts as a lookup between the
// names of our custom template functions and the functions themselves.
// Parse the base template page into a template set.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

// Function that returns a cache containing html templates and a customer
// FuncMap for use in generating response output.
func NewTemplateCache() (map[string]*template.Template, error) {
	// Initialise a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use the fs.Glob() to get a slice of all filepaths in the ui.Files embedded
	// filesystem which match the pattern "html/pages/*.tmpl".  This essentially gives
	// us a slice of all the page templates just like before.
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	// Loop through the page filepaths one-by-one.
	for _, page := range pages {
		// Extract the file name (like 'home.tmpl') from the full filepath
		// and assign it to the name variable.
		name := filepath.Base(page)

		// Create a slice containing the filepath patterns for the templates we
		// want to parse.
		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		// Use ParseFS() instead of ParseFiles() to parse the template files
		// from the ui.Files embedded filesystem.
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map, using the name of the page
		// (like 'home.tmpl') as the key.
		cache[name] = ts
	}

	// Return the map.
	return cache, nil
}
