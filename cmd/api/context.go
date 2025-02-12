package main

import (
	"SrbastianM/rest-api-gin/internal/data"
	"context"
	"net/http"
)

// Define a custom contextKey type, with the underlying type string.
type contextKey string

// Conevrt the string "user" to a contextKey type and assing it to the userContextKey constant.
const userContextKey = contextKey("user")

// Returning a new copy of the request with the provided User struct added to the context.
func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// Retrieves the User struct from the request context.
func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
