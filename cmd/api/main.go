package main

import (
	"SrbastianM/rest-api-gin/internal/data"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// declare the version for the api -> 1.0.0 first version
const version = "1.0.0"

// define the struct to hold the config settings:
// server port: where the app listen on,
// current environment: stage, production, etc...
type config struct {
	port int
	env  string
	db   struct {
		dsn           string
		maxOpensConns int
		maxIdleConns  int
		maxIdleTime   string
	}
}

// Define the struct to hold the dependencies for the HTTP handlers,
// helpers and middleware
// Contain: copy of the config struct and a logger (just for now)
type application struct {
	config config
	logger *log.Logger
	models data.Models
}

func main() {
	// instance of the config struct
	var cfg config
	// Read the value of the port and env command-flag lines into config struct
	// For default i use port: 4000 and the environment "development"
	flag.IntVar(&cfg.port, "port", 4000, "API Server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	// Read the DSN value from the db-dsn command-line flag into the config struct. It use development DSN cuz
	//  is not flag provided
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading the .env file: %s", err)
	}
	flag.StringVar(&cfg.db.dsn, "db-sn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")

	flag.IntVar(&cfg.db.maxOpensConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Parse()

	//initialize new logger which writes messages to the standard out stream
	// prefixed with the current date and time
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	logger.Printf("database connection pool established")
	// instance of the aplication struct, contains config struct and the logger
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	// Declare a new serveMux and add the version and the route wich dispatches requests
	// to healthcheck method
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)

	// Declare a HTTP server with some sensible properties timeout settings
	// it listents on the port provided (4000), and use the httprouter as handler
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start the HTTP Server
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)

}

// Return the sql.DB connection pool
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	// Create an empty connection pool, using the DSN from the config
	if err != nil {
		return nil, err
	}
	// Set the maximum number of open (in+use + idle) connection in the pool.
	db.SetMaxOpenConns(cfg.db.maxOpensConns)
	// Set the maximun number of idle connections in the pool.
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	// Use the time.parseDuration() function to convert the idle timeout duration string to
	// time.Duration type.
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	// Create a context with 5 seconds timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Establish a new connection to the database, passing in the context
	// created above as a parameter. If couldnt be established succesfully
	// return an error (5 seconds deadline)
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	// Return the sql.DB connection pool.
	return db, nil
}
