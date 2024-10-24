package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Register the methods, URL patterns and handler functions
// for the end points GET and POST using the HandlerFunction() method
// and return the httprouter instance
func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/foods", app.createFoodHandler)
	router.HandlerFunc(http.MethodGet, "/v1/foods/:id", app.showFoodHandler)

	return router
}
