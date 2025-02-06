package main

import (
	"fmt"
	"net/http"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	// Create a defeared function wich always be run in the event of a panic as Go
	// unwinds the stack.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// Use the builtin recover function to check if there has been a panic or not
			if err := recover(); err != nil {
				// If there a panic, set a "Connection: close" header on the response. This
				// act as a trigger to mage Go's HTTP server automatically close the current connection
				// afeter a response has been sent.
				w.Header().Set("Connection", "close")
				// The value returned by recover() has the type interface{}, so we use fmt.Error() to
				// normalize it into an error and call our serverErrorResponse() helper
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
