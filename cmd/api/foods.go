package main

import (
	"fmt"
	"net/http"
)

func (app *application) createFoodHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create new food")
}

// Use the helper "ReadIdParam"
func (app *application) showFoodHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "show food details %d\n", id)
}
