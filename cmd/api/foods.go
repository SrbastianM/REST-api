package main

import (
	"SrbastianM/rest-api-gin/internal/data"
	"fmt"
	"net/http"
	"time"
)

func (app *application) createFoodHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create new food")
}

// Use the helper "ReadIdParam"
func (app *application) showFoodHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Create a new instance of the MOvie Struct, containing the ID extracted from
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
