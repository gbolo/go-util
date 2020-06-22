package backend

import (
	"fmt"
	"net/http"

	_ "appname/docs"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title appname
// @version 0.1
// @description Swagger API appname
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email dev@appname
// @BasePath /api

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

var routes = []Route{
	{
		"Version",
		"GET",
		getEndpoint("version"),
		handlerVersion,
	},
}

func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {

		// add compression support to handler if enabled
		var handler http.Handler
		handler = route.HandlerFunc
		if viper.GetBool("server.compression") {
			handler = handlers.CompressHandler(route.HandlerFunc)
		}

		// add routes to mux
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	// add swagger UI
	router.Methods("GET").Path("/swagger").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// redirect to /swagger/index.html which is provided by httpSwagger.WrapHandler
		http.Redirect(w, req, "/swagger/index.html", 301)
	})
	router.Methods("GET").PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	// add route to mux to handle frontend UI static files (generated by npm)
	staticPath := viper.GetString("server.static_files_dir")
	if staticPath == "" {
		staticPath = "./frondent/dist"
	}

	handlerStatic := http.StripPrefix("/", http.FileServer(http.Dir(staticPath)))
	// add compression support to handler if enabled
	if viper.GetBool("server.compression") {
		handlerStatic = handlers.CompressHandler(handlerStatic)
	}

	router.
		Methods("GET").
		PathPrefix("/").
		Handler(handlerStatic)

	return router
}