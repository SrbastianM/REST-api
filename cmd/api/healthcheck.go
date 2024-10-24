package main

import (
	"fmt"
	"net/http"
)

// Show the application information
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Create a fixed JSON response from a string
	// wrap the interpolated values in double quotes in to %q verb
	js := `{"status": "available", "environment": %q, "version": %q}`
	js = fmt.Sprintf(js, app.config.env, version)

	//Set the default sending a "Content-Type: text/plain; charset=utf-8"
	w.Header().Set("Content-Type", "application/json")

	// Write JSON as the HTTP response body.
	w.Write([]byte(js))
}
