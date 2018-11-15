package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// returns all envs
func handlerEnv(w http.ResponseWriter, req *http.Request) {

	jData, err := json.Marshal(cachedTable)
	if err != nil {
		log.Errorf("error parsing json: %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

// handler to powerdown an env
func handlerEnvPowerdown(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// get vars from request to determine environment
	vars := mux.Vars(req)
	envName := vars["env"]

	if envName != "" {
		res, err := shutdownEnv(envName)
		writeJsonResponse(w, err, res)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"error\":\"empty environment name\"}\n")
	}
}

// handler to start up an env
func handlerEnvStartup(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// get vars from request to determine environment
	vars := mux.Vars(req)
	envName := vars["env"]
	log.Infof("starting env: %s", envName)

	if envName != "" {
		res, err := startupEnv(envName)
		writeJsonResponse(w, err, res)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"error\":\"empty environment name\"}\n")
	}
}

// handler to refresh envs
func handlerEnvRefresh(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := refreshTable(); err != nil {
		log.Errorf("refresh error: %v", err)
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
	} else {
		log.Info("refresh successful")
		fmt.Fprint(w, "{\"status\":\"OK\"}\n")
	}
}

func writeJsonResponse(w http.ResponseWriter, err error, response []byte) {
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		if len(response) > 0 {
			w.Write(response)
		} else {
			fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		}
	}
}
