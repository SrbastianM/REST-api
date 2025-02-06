package main

import (
	"fmt"
	"net/http"
)

// the logError() method is a generic logging error message.
func (app *application) logError(r *http.Request, err error) {
	app.logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

// the errorResponse() method is a generic helper for sending JSON-formatted error
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelop{"error": message}

	// write a error if this happens and return 500 internal server status code
	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

// the serverErrorResponse() method is a generic helper for sending a message if the server
// encountered an unexpected problem at runtime. Return a message and 500 internal server error
// status code. The message is a JSON response that contain the generic error
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "The server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// the serverErrorResponse() method will used to send 404 status code and
// JSON response to the client
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {

	message := "the requested resource could not found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// the methodNotAllowedResponse() method will used to send 405 Method not allowed
// and JSON response to the client
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {

	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

// the method editConflictResponse() will used to send and 409 Conflict response.
func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusConflict, message)
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	app.errorResponse(w, r, http.StatusTooManyRequests, message)
}
