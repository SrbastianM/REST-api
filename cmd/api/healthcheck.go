package main

import (
	"net/http"
)

// Show the application information
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Create a map which hold the information that want to send in the response
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	// implement the helper writeJson to send the "Content-type: application-json" header
	// and send the data encode to JSON and response the status and the JSON Response -> see the writeJSON() helper
	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
