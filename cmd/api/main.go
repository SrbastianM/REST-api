package main

import (
	"SrbastianM/rest-api-gin/internal/data"
	"SrbastianM/rest-api-gin/internal/jsonlog"
	"SrbastianM/rest-api-gin/internal/mailer"
	"context"
	"database/sql"
	"flag"
	"log"
	"os"
	"sync"
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
	limiter struct {
		rps    float64
		burst  int
		enable bool
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

// Define the struct to hold the dependencies for the HTTP handlers,
// helpers and middleware
// Contain: copy of the config struct and a logger (just for now)
type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
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
		log.Fatalf("Error loading the .env file: %s", err)
	}
	flag.StringVar(&cfg.db.dsn, "db-sn", os.Getenv("DB_DSN"), "PostgreSQL DSN")

	flag.IntVar(&cfg.db.maxOpensConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximun request per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximun burst")
	flag.BoolVar(&cfg.limiter.enable, "limiter-enable", true, "Enable rate limiter")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "f16468f38c5882", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "da68096b0dc8fd", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Foody <no-reply@foody.net>", "SMTP sender")

	flag.Parse()

	//initialize new logger which writes messages to the standard out stream
	// prefixed with the current date and time
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()

	logger.PrintInfo("database connection pool established", nil)
	// instance of the aplication struct, contains config struct and the logger
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
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
