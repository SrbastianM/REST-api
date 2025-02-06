package main

import (
	"net/http"
)

// Show the application information
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Declare a envelop map containing the data for the response.
	// Create a map which hold the information that want to send in the response like a nested JSON
	env := envelop{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	// implement the helper writeJson to send the "Content-type: application-json" header
	// and send the data encode to JSON and response the status and the JSON Response -> see the writeJSON() helper
	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.logger.PrintError(err, nil)
		app.serverErrorResponse(w, r, err)
	}
}
