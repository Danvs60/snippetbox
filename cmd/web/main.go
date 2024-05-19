package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/Danvs60/snippetbox/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

// Define an application struct to hold application-wide dependencies
// This is for dependency injection to avoid using global variables
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
	templateCache map[string]*template.Template
}


func main() {
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	dsn := flag.String("dsn", "web:st0ngb00ze@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	// flags and local date and local time (joined by the bitwise OR |)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// include Lshortfile to include file name and line number of error
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()


	// Template cache...
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}
	// Initialise application var, the dependency container
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	// Use ListenAndServe() to start a new web server.
	// Takes in 2 parameters:
	//	1. the TCP network address to listen on
	//	2. and the servemux (mapper of URLs) (actually it accepts a Handler [interface that implements ServeHTTP]
	//		as ServeMux does actually implement a ServeHTTP function, then it is fine to use as parameter
	//	think of mux as a handler for handlers
	// If mux returns an error we log it.

	// err := http.ListenAndServe(*addr, mux)

	// initialise new http.Server to use custom logger
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	app.infoLog.Printf("Starting server on %s", srv.Addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

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
