package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	// Create a defeared function wich always be run in the event of a panic as Go
	// unwinds the stack.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// Use the builtin recover function to check if there has been a panic or not
			if err := recover(); err != nil {
				// If there a panic, set a "Connection: close" header on the response. This
				// act as a trigger to mage Go's HTTP server automatically close the current connection
				// afeter a response has been sent.
				w.Header().Set("Connection", "close")
				// The value returned by recover() has the type interface{}, so we use fmt.Error() to
				// normalize it into an error and call our serverErrorResponse() helper
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	// Initialize a new rate limiter wich allows an average of 2 request per second,
	// with a maximun of a 4 request in a single burst -> at the same time
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the clients IP address from the request
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		// Lock the mutex to prevent this code form being execuded concurrently.
		mu.Lock()
		// Check if the IP is already on the map, if doesnt then initialize a new rate
		// limiter and add the IP and the limiter to the map
		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: rate.NewLimiter(2, 4)}
		}

		clients[ip].lastSeen = time.Now()
		// Check if the IP address exists. If the request isn't allowed, ulock the mutext
		// and return the method rateLimitExceededResponse() wich means that is to many request
		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			app.rateLimitExceededResponse(w, r)
			return
		}

		mu.Unlock()
		next.ServeHTTP(w, r)
	})
}
