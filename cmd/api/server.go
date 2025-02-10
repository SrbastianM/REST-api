package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)
	// Start a goroutine wich send a notify to listen the incoming SIGNINT and SIGTERM signal
	// then relay them to the quit channel. Return a message to say the graful shutdown will initialize.
	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit
		app.logger.PrintInfo("shutting down server", map[string]string{
			"signal": s.String(),
		})
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	app.logger.PrintInfo("Starting serve", map[string]string{
		"add": srv.Addr,
		"env": app.config.env,
	})
	// Safety check for the Shutdown() method.
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// If the shutdownError chan has any problem, it catch and returning what is happening.
	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.PrintInfo("stopped serve", map[string]string{
		"addr": srv.Addr,
	})
	return nil
}
