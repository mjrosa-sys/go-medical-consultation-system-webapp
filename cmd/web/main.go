package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type flags struct {
	port string
	dsn  string
}

type application struct {
	flags flags
	db    *sql.DB
}

func main() {
	app := application{}

	flag.StringVar(&app.flags.port, "port", ":3001", "Port where the application is gonna run.")
	flag.StringVar(&app.flags.dsn, "dsn", "gocare_app:123@tcp(localhost:3306)/gocare-db", "MySQL data source name")
	flag.Parse()

	db, err := openDB(app.flags.dsn)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	app.db = db

	srv := &http.Server{
		Addr:         app.flags.port,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Server running localhost%s\n", app.flags.port)

	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	if err != nil {
		log.Fatal(err.Error())
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
