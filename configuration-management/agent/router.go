package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	// APIVersion defines the compatability version of the API and is appended to each API route
	APIVersion     = "1"
	endpointFormat = "/api/v%s/%s"
)

// getEndpoint returns a properly formatted API endpoint
func getEndpoint(suffix string) string {
	return fmt.Sprintf(endpointFormat, APIVersion, suffix)
}

// Route defines a route passed to our mux
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes holds a list of Routes
type Routes []Route

// all defined server endpoints
var routes = Routes{

	// API endpoints
	Route{
		"Version",
		"GET",
		getEndpoint("version"),
		handlerVersion,
	},

	// HQ will POST tasks for the agent on this endpoint
	Route{
		"DoTask",
		"POST",
		getEndpoint("task"),
		handlerTask,
	},
}

func newRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = handlers.CompressHandler(route.HandlerFunc)

		// add routes to mux
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
