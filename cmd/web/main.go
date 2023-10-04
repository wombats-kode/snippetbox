package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"snippetbox/internal/models"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql" // Manual MySQL import
)

// Add a formDecoder field to hold a pointer to a form.Decoder instance

// Define an application struct to hold the application-wide dependencies for the
// Web application.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	// Use our Interface type so we can leverage mocks in testing
	snippets       models.SnippetModelInterface
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pa55w0rd@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	// Define custom error logger to differentiate INFO and ERROR logging.  Set one to
	// log to Stderr and the other to log to stdout so they can be differentiated into
	// separate file like:
	// go run ./cmd/web >>/tmp/info.log 2>>/tmp/error.log
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Create an object to hold the MySQL connection pool
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	// We defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.
	defer db.Close()

	// Initialise a new templateCache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}
	formDecoder := form.NewDecoder()

	// Use the scs.New() function to initialise a new session manager. Then we
	// configure it to use our MySQL database as the session store, and set a
	// lifetime of 12 hours (so that the sessions automatically expire 12 hours
	// after first being created).
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// Make sure that the Secure attribute is set on our session cookies
	// Setting this means that the cookie will only  be sent by a user's web
	// browser when a HTTPS connection that is being used
	sessionManager.Cookie.Secure = true

	// and add to the app dependencies

	app := application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// Initialise a tls.Config struct to hold the non-default TLS settings we
	// want the server to use.  Here we are optimiseing the struct to use onley elliptic
	// curves that with assembly implementations, for performance optimisation.
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS12,
		// CipherSuite settings will only impact TLS1.2 or lower.
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	// Initialise a new http.server struct.  We set the addr and Handler fields so
	// that the server uses the same network address and routes as before, and set the
	// Errorlog field so the server now uses the custom errorLog logger in the event
	// of any problems
	srv := &http.Server{
		Addr:      *addr,
		ErrorLog:  errorLog,
		Handler:   app.routes(),
		TLSConfig: tlsConfig,
		// Add Idle, Read and Write timeouts to the server
		IdleTimeout:    time.Minute,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 524288, // default is 1MB
	}

	app.infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

// The openDb() function wraps sql.open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
