package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}


// returns version information
func handlerVersion(w http.ResponseWriter, req *http.Request) {
	writeJSONResponse(w, http.StatusOK, map[string]string{"version": "alpha"})
}

// executes a task
func handlerTask(w http.ResponseWriter, req *http.Request) {
	// try to read the body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		apiResponse := errorResponse{"Bad request. Cannot read request body."}
		writeJSONResponse(w, http.StatusBadRequest, apiResponse)
		return
	}

	// try to unmarshal the body into a task
	task := task{}
	if err = json.Unmarshal(body, &task); err != nil {
		apiResponse := errorResponse{"Bad request. Cannot decode request body."}
		writeJSONResponse(w, http.StatusBadRequest, apiResponse)
		return
	}

	switch task.Module {
	// these are the only supported modules
	case "directory":
		writeJSONResponse(w, http.StatusOK, ModuleDirectory(task))
	case "file":
		writeJSONResponse(w, http.StatusOK, ModuleFile(task))
	case "apt":
		writeJSONResponse(w, http.StatusOK, ModuleApt(task))
	case "service":
		writeJSONResponse(w, http.StatusOK, ModuleService(task))
	// we dont support anything else but the above
	default:
		log.Errorf("unrecognized module specified: %s", task.Module)
		writeJSONResponse(w, http.StatusBadRequest, errorResponse{"invalid module"})
	}
}

// wrapper for json responses
func writeJSONResponse(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	b, _ := json.MarshalIndent(body, "", "  ")
	w.Write(b)
}
