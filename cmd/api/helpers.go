package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Retrieve the "id" parameter form the current request context then convert it into a integer
// and return it. If the operation isn't successfully, return 0 an error.
func (app *application) readIDParam(r *http.Request) (int64, error) {

	// parsing a request, any interpolate URL parameters will stored in the request context
	// and retrieve an slice that contain the names and values
	params := httprouter.ParamsFromContext(r.Context())

	// Use the ByName() method to get the value of the "id" parameter from slice
	// Convert the integer in to a base integer (10 integer with a bit size of 64)
	// Check if the value couldn't not convert or is less than 1
	// and handle the error using http.NotFound() and return 404 response
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid parameter")
	}

	return id, nil
}
