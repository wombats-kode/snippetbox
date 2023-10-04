package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"snippetbox/internal/models"
	"snippetbox/ui"
)

// Define a templateData type to act as the holding structure or
// any dynamic data that we want to pass to our HTML templates

type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string // Add a Flash field to the templateData struct
	IsAuthenticated bool
	CSRFToken       string // adds a CSRFToken field
}

// Create a humanDate() function which returns a nicely formatted string.
func humanDate(t time.Time) string {
	// Return the empty string if time has the zero value.
	if t.IsZero() {
		return ""
	}

	// Convert the time to UTC before formatting it
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

// Initialise a template.FuncMap object and store it in a global variable.
// This is essentially a string-keyed map which acts as a lookup between the
// names of our custom template functions and the functions themselves
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialise a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use fs.Glob() to get a slice of all the filepaths in the ui.Files embedded
	// filesystem which match the pattern 'html/pages/*.tmpl'.  This essentially
	// gives us a slice of all the 'page' templates for the application, just like before.
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	// Loop through the page filepaths one-by-one.
	for _, page := range pages {
		// Extract the file name (like 'home.tmpl') from the fullpath
		// and assign it to the name variable
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
		// Like 'home.tmpl' as the key.
		cache[name] = ts
	}
	return cache, nil
}
