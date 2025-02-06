package main

import (
	"SrbastianM/rest-api-gin/internal/data"
	"SrbastianM/rest-api-gin/internal/validator"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Food struct {
	ID        int64
	CreatedAt time.Time
	Title     string
	Types     []string
	Version   int32
}

func (app *application) createFoodHandler(w http.ResponseWriter, r *http.Request) {
	// Declare a anonyms struct to hold the information that expect to be in HTTP
	// request body
	var input struct {
		Title string
		Types []string
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

	food := &data.Food{
		Title: input.Title,
		Types: input.Types,
	}

	v := validator.New()

	if data.ValidateFood(v, food); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Call Insert() method in the food model, passing pointer to the validated movie struct.
	// This wil create a record in the database and update the movie struct with the
	// system-generated information.
	err = app.models.Foods.Insert(food)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Make an empty http.Header map and then use the Set() method to add a new location Header,
	// interpoling the system-generated ID for the new food in the URL
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v/foods/%d", food.ID))

	// Wrte a Json response with a 201 Created status code, the food data in the
	// response body, and the location header
	err = app.writeJSON(w, http.StatusCreated, envelop{"food": food}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showFoodHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Change the hardcode value to the Get() method create in internal/data/foods.go.
	// This method fetch the data for a specific food. Also catch some error if is not found,
	// returning 404 Not Found response to the client.
	food, err := app.models.Foods.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// encode the struct to JSON  and send it as the HTTP Response
	// Create an instance of envelop to it to writeJSON()
	err = app.writeJSON(w, http.StatusOK, envelop{"food": food}, nil)
	if err != nil {
		app.logger.PrintError(err, nil)
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateFoodHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	food, err := app.models.Foods.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if r.Header.Get("X-expected-Version") != "" {
		if strconv.FormatInt(int64(food.Version), 32) != r.Header.Get("X-Expected-Version") {
			app.editConflictResponse(w, r)
		}
		return
	}

	// Declare a struct to hold the expected data from client
	var input struct {
		Title *string
		Types []string
	}
	// Read the json request body data into the input struct
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// If the input value is nil it means that no corresponding "title" key/value pair provided
	// by the json reques body. So  move on and leave the record unchanged. Otherwise, update the
	// record with the new title value. It happens too with the types slices.
	if input.Title != nil {
		food.Title = *input.Title
	}

	if input.Types != nil {
		food.Types = input.Types
	}

	// Validate the updated food record, sending the client a 422 Unprocessable Entity
	// response if any checks fails
	v := validator.New()

	if data.ValidateFood(v, food); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Pass the updated movie record to our new Update() method. Intercept any ErrEditConflict() error and
	// call the new editConflictResponse() helper.
	err = app.models.Foods.Update(food)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Write the updated movie record in a JSON Response
	err = app.writeJSON(w, http.StatusOK, envelop{"food": food}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteFoodHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the food Id from the URL
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Delete the food from DB, sending 404 Not found response to the client if there isn't a
	// matching record
	err = app.models.Foods.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Return 200 OK status code along with a success message
	err = app.writeJSON(w, http.StatusOK, envelop{"message": "Food successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listFoodHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string
		Types []string
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()
	// Use the helpers to extract the title, the types string values, falling back to default of an
	// empty string and an empty slice respectively it they are not provided by the client
	input.Title = app.readString(qs, "title", "")
	input.Types = app.readCSV(qs, "types", []string{})
	// Get the page size query string values as integers and set the default page value to 1 and
	// the default size to 20.
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "title", "-id", "-title", "type"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	foods, metadata, err := app.models.Foods.GetAll(input.Title, input.Types, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelop{"foods": foods, "metada": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
