package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/kamil-alekber/deliciousness-recipes/internal/models/recipes"
	_ "modernc.org/sqlite"
)

type config struct {
	dsn       string
	addr      string
	staticDir string
}

var cfg config

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	recipes       *recipes.Queries
	templateCache map[string]*template.Template
}

func initDbTables(db *sql.DB) error {
	// return all folder names inside internal/models/
	dir, err := os.ReadDir("internal/sql")
	if err != nil {
		return err
	}

	for _, entry := range dir {
		if entry.IsDir() {
			path := fmt.Sprintf("internal/sql/%s/schema.sql", entry.Name())
			schema, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			ddl := string(schema)

			if _, err := db.ExecContext(context.Background(), ddl); err != nil {
				return err
			}

			log.Printf("Executed scema: %s", entry.Name())
		}
	}

	return nil
}

func main() {
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&cfg.dsn, "dsn", "store.db", "mysql data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(cfg.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	if err := initDbTables(db); err != nil {
		errorLog.Fatal(err)
	}

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		recipes:       recipes.New(db),
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", cfg.addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
