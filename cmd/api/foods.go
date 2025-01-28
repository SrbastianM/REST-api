package main

import (
	"SrbastianM/rest-api-gin/internal/data"
	"SrbastianM/rest-api-gin/internal/validator"
	"errors"
	"fmt"
	"net/http"
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

// Use the helper "ReadIdParam"
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
		app.logger.Println(err)
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
	// Declare a struct to hold the expected data from client
	var input struct {
		Title string
		Types []string
	}
	// Read the json request body data into the input struct
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Copy the values from the request body to the aproppiate fields of the food
	food.Title = input.Title
	food.Types = input.Types

	// Validate the updated food record, sending the client a 422 Unprocessable Entity
	// response if any checks fails
	v := validator.New()

	if data.ValidateFood(v, food); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Pass the updated movie record to our new Update() method
	err = app.models.Foods.Update(food)
	if err != nil {
		app.serverErrorResponse(w, r, err)
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
