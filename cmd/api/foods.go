package main

import (
	"SrbastianM/rest-api-gin/internal/data"
	"fmt"
	"net/http"
	"time"
)

func (app *application) createFoodHandler(w http.ResponseWriter, r *http.Request) {
	// Declare a anonyms struct to hold the information that expect to be in HTTP
	// request body
	var input struct {
		Title  string
		Year   int32
		Genres []string
	}

	// Initialize a new json.Decoder() instance which reads from the request body
	// and then use Decode() method to decode the body contents into input struct.
	// If there was an error during decoding it send a generic errorResponse() helper
	// to send the 400 bad request
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

// Use the helper "ReadIdParam"
func (app *application) showFoodHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Create a new instance of the food Struct, containing the ID extracted from
	// URL and dummy data.
	food := data.Food{
		ID:       id,
		CreateAt: time.Now(),
		Title:    "Potato",
		Types:    []string{"vegetables", "fruit", "fat"},
		Version:  1,
	}
	// encode the struct to JSON  and send it as the HTTP Response
	// Create an instance of envelop to it to writeJSON()
	err = app.writeJSON(w, http.StatusOK, envelop{"food": food}, nil)
	if err != nil {
		app.logger.Println(err)
		app.serverErrorResponse(w, r, err)
	}
}
