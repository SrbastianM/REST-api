package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// declare the version for the api -> 1.0.0 first version
const version = "1.0.0"

// define the struct to hold the config settings:
// server port: where the app listen on,
// current environment: stage, production, etc...
type config struct {
	port int
	env  string
}

// Define the struct to hold the dependencies for the HTTP handlers,
// helpers and middleware
// Contain: copy of the config struct and a logger (just for now)
type application struct {
	config config
	logger *log.Logger
}

func main() {
	// instance of the config struct
	var cfg config
	// Read the value of the port and env command-flag lines into config struct
	// For default i use port: 4000 and the environment "development"
	flag.IntVar(&cfg.port, "port", 4000, "API Server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	//initialize new logger which writes messages to the standard out stream
	// prefixed with the current date and time
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// instance of the application struct, contain config struct and the logger
	app := &application{
		config: cfg,
		logger: logger,
	}
	// Declare a new serveMux and add the version and the route wich dispatches requests
	// to healthcheck method
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)

	// Declare a HTTP server with some sensible properties timeout settings
	// it listents on the port provided (4000), and use the servemux as handler
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start the HTTP Server
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)

}
