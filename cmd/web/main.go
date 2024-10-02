package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"snippetbox.example.com/internal/models"
)

// Define an application struct to hold the application-wide dependencies for the
// web application.
// Holds an slog logger dependency
// Hold a models snippetmodel dependency
// Holds a ui templateCache
type application struct {
	logger        *slog.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:w3bpassword@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// To keep the main() function tidy we will use a separate openDB() function
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	// Initialise a new template cache...
	templateCache, err := NewTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		logger:        logger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder:   form.NewDecoder(),
	}

	logger.Info("starting server", "addr", *addr)

	// Call the new app.Routes() method to the servemux containing our routes,
	// and pass that to http.ListenAndServe().
	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

// The openDB function wraps sql.Open and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
