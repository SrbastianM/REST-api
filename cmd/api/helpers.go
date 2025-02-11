package main

import (
	"SrbastianM/rest-api-gin/internal/validator"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// use http.MaxBytesReader() to limit the size of the request body to 1MB
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// if there is an error during the decoding, start the triage
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		// use the errors.AS() function whether the error has the type *json.SyntaxError
		// if it does, return a plain text error message which includes the location of the problem
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
			// If the Decode() return an io.ErrUnexpectedEOF error for syntax errors in the JSON
			// it returns a generic error message
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
			// If the Decode() return an json.UnmarshalTypeError error for JSON value is the wrong type for the
			// target destination
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type (at character %q)", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
			// If the request body is empty return a io.EOF error
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		// If the JSON contains a field wich cannot be mapped to the target destination
		//	then Decode() will now return an error message
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		// If the request exceeds 1MB in size the decode will now fail
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger dan %d bytes", maxBytes)
			// If pass a no nil pointer to Decode() catch and panic
			// rather returning the error to the handler
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	return nil
}

// This helper returns a string value from the query string, or the provided default value if
// no matching key could be found.
func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}
	return s
}

// This helper reads a stirng from the query string and then splits it into a slice on the comma character.
// If no matching key could be found, it returns the provided default value
func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)
	if csv == "" {
		return defaultValue
	}
	return strings.Split(csv, "")
}

// This helper reads an string value from the query string and converts it to an integer before returning.
// If no matching key could be found it returns the provided default value. If couldnt be converter to an integer
// then record an error message in the provided Validator instance
func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "Must be an integer value")
		return defaultValue
	}

	return i
}

// Accepts an arbitrary function as a parameter and recover the error instead of panic the aplication
func (app *application) background(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				app.logger.PrintError(fmt.Errorf("%s", err), nil)
			}
		}()
		fn()
	}()

}
