package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Define the envelop type
type envelop map[string]interface{}

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

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelop, header http.Header) error {
	//Encode the data to JSON and return err if there was one
	// Indent the terminal output putting a tab and no line prefix ("")
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	// We won't encored errors at this point so we can write any headers.
	// we loop through the header map and add each header to the http.ResponseWriter header map.
	// if the map isn't nil that means its OK. Go doesn't trow an error
	for key, value := range header {
		w.Header()[key] = value
	}
	// Add "Content-Type: application-json" header, them write the status code and JSON response
	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
