package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Register the methods, URL patterns and handler functions
// for the end points GET and POST using the HandlerFunction() method
// and return the httprouter instance
func (app *application) routes() http.Handler {
	router := httprouter.New()

	// convert the method helper to a http.Handler and set the custom
	// handler  error (404) not found
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	// likewise, convert the method methodNotAllowedResponse helper to a http.Handler
	// and set the custom handler (405) not allowed responses
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/foods", app.createFoodHandler)
	router.HandlerFunc(http.MethodGet, "/v1/foods", app.listFoodHandler)
	router.HandlerFunc(http.MethodGet, "/v1/foods/:id", app.showFoodHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/foods/:id", app.updateFoodHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/foods/:id", app.deleteFoodHandler)

	return app.recoverPanic(app.rateLimit(router))
}
