// +build go1.8

// enforce go 1.8+ just so we can support X25519 curve :)

package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

var (
	// set timeouts to avoid Slowloris attacks.
	httpWriteTimeout = time.Second * 60
	httpReadTimeout  = time.Second * 15
	// the maximum amount of time to wait for the
	// next request when keep-alives are enabled
	httpIdleTimeout = time.Second * 60
)

func startHTTPServer() (err error) {

	// create routes
	mux := newRouter()

	// get server config
	srv := configureHTTPServer(mux)

	// get TLS config
	tlsConifig, err := configureTLS()
	if err != nil {
		log.Fatalf("error configuring TLS: %s", err)
		return
	}
	srv.TLSConfig = &tlsConifig

	// start the server
	if viper.GetBool("server.tls.enabled") {
		// cert and key should already be configured
		log.Infof("starting HTTP server with TLS enabled: listening on %s", srv.Addr)
		err = srv.ListenAndServeTLS("", "")
	} else {
		log.Infof("starting HTTP server: listening on %s", srv.Addr)
		err = srv.ListenAndServe()
	}

	if err != nil {
		log.Fatalf("failed to start server: %s", err)
	}

	return
}

func configureHTTPServer(mux *mux.Router) (httpServer *http.Server) {

	// apply standard http server settings
	address := fmt.Sprintf(
		"%s:%s",
		viper.GetString("server.bind_address"),
		viper.GetString("server.bind_port"),
	)

	httpServer = &http.Server{
		Addr: address,

		WriteTimeout: httpWriteTimeout,
		ReadTimeout:  httpReadTimeout,
		IdleTimeout:  httpIdleTimeout,
	}

	// explicitly enable keep-alives
	httpServer.SetKeepAlivesEnabled(true)

	// stdout access log enable/disable
	if viper.GetBool("server.access_log") {
		httpServer.Handler = handlers.CombinedLoggingHandler(os.Stdout, mux)
	} else {
		httpServer.Handler = mux
	}

	return
}