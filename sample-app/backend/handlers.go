package backend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/mux"
)

// @Summary Returns version information
// @Description Returns version information
// @Tags Misc
// @Produce json
// @Success 200 {object} versionInfo
// @Router /v1/version [get]
func handlerVersion(w http.ResponseWriter, req *http.Request) {
	writeJSONResponse(w, http.StatusOK, getVersionResponse())
}

// wrapper for json responses
func writeJSONResponse(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	b, _ := json.MarshalIndent(body, "", "  ")
	w.Write(b)
}

// @Summary Returns the health of the application
// @Description If the status is not 200, then the application is unhealthy
// @Tags Misc
// @Produce json
// @Success 200 {object} successResponse "server is healthy"
// @Failure 500 {object} errorResponse "server is unhealthy. Usually a database issue"
// @Router /v1/healthz [get]
func handlerHealthCheck(w http.ResponseWriter, req *http.Request) {
	if err := dbHealthCheck(); err != nil {
		message := fmt.Sprintf("error connecting to database: %v", err)
		writeJSONResponse(w, http.StatusInternalServerError, errorResponse{message})
		return
	}
	writeJSONResponse(w, http.StatusOK, successResponse{"OK"})
}

// @Summary Add a new Client
// @Description Add a new Client
// @Tags Clients
// @Produce json
// @Param client body Client true "Add Client"
// @Success 201 {object} successResponse "new client was added"
// @Failure 409 {object} errorResponse "a client with that ID already exists. Try running an update instead"
// @Failure 400 {object} errorResponse "invalid request body"
// @Failure 500 {object} errorResponse "server was unable to process the request. Usually a database issue"
// @Router /v1/client [post]
func handlerAddClient(w http.ResponseWriter, req *http.Request) {
	// read in body
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Errorf("error reading request body: %s", err)
		writeJSONResponse(w, http.StatusBadRequest, errorResponse{"request body could not be read"})
		return
	}
	req.Body.Close()

	// decode json body
	var body Client
	err = json.Unmarshal(reqBody, &body)
	if err != nil {
		log.Errorf("failed to decode body: %s", err)
		writeJSONResponse(w, http.StatusBadRequest, errorResponse{"request body could not be decoded"})
		return
	}

	// validate body
	_, err = govalidator.ValidateStruct(body)
	if err != nil {
		log.Debugf("failed to validate body: %s", err)
		writeJSONResponse(w, http.StatusBadRequest, errorResponse{"request body failed input validation"})
		return
	}

	// database transactions...
	// check if client already exists
	var existingClient *Client
	existingClient, err = dbGetClient(body.ID)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}
	if existingClient != nil && existingClient.ID == body.ID {
		writeJSONResponse(w, http.StatusConflict, errorResponse{"client with specified ID already exists"})
		return
	}
	// this is a new client, so create it
	err = dbCreateClient(&body)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}
	writeJSONResponse(w, http.StatusCreated, successResponse{"OK"})
}

// @Summary Updates an existing Client
// @Description Updates an existing Client with the specified ID
// @Tags Clients
// @Produce json
// @Param client body Client true "Add Client"
// @Success 200 {object} successResponse "client has been updated"
// @Success 304 {object} successResponse "client does not need updating, it's already at that state"
// @Failure 404 {object} errorResponse "a client with that ID does not exist"
// @Failure 400 {object} errorResponse "invalid request body"
// @Failure 500 {object} errorResponse "server was unable to process the request. Usually a database issue"
// @Router /v1/client [put]
func handlerUpdateClient(w http.ResponseWriter, req *http.Request) {
	// read in body
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Errorf("error reading request body: %s", err)
		writeJSONResponse(w, http.StatusBadRequest, errorResponse{"request body could not be read"})
		return
	}
	req.Body.Close()

	// decode json body
	var body Client
	err = json.Unmarshal(reqBody, &body)
	if err != nil {
		log.Errorf("failed to decode body: %s", err)
		writeJSONResponse(w, http.StatusBadRequest, errorResponse{"request body could not be decoded"})
		return
	}

	// validate body
	_, err = govalidator.ValidateStruct(body)
	if err != nil {
		log.Debugf("failed to validate body: %s", err)
		writeJSONResponse(w, http.StatusBadRequest, errorResponse{"request body failed input validation"})
		return
	}

	// database transactions...
	// ensure that the client already exists
	var existingClient *Client
	existingClient, err = dbGetClient(body.ID)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}
	if existingClient == nil || existingClient.ID == "" {
		writeJSONResponse(w, http.StatusNotFound, errorResponse{"client with specified ID does not exist"})
		return
	}
	if cmp.Equal(existingClient, body) {
		writeJSONResponse(w, http.StatusNotModified, errorResponse{"client already exists with specified parameters, no updated needed"})
		return
	}
	// client needs updating
	err = dbUpdateClient(&body)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}
	writeJSONResponse(w, http.StatusOK, successResponse{"OK"})
}

// @Summary Returns a list of all clients
// @Description Returns a list of all clients
// @Tags Clients
// @Produce json
// @Success 200 {array} Client "a list of clients"
// @Failure 500 {object} errorResponse "an error occurred. Usually a database issue"
// @Router /v1/client [get]
func handlerGetClients(w http.ResponseWriter, req *http.Request) {
	// database transaction
	body, err := dbGetAllClients()
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}
	writeJSONResponse(w, http.StatusOK, body)
}

// @Summary Delete a Client by ID
// @Description Delete a Client by ID
// @Tags Clients
// @Produce json
// @Param id path string true "Client ID"
// @Success 200 {object} successResponse "client was deleted or does not exist"
// @Failure 500 {object} errorResponse "an error occurred. Usually a database issue"
// @Router /v1/client/{id} [delete]
func handlerDeleteClient(w http.ResponseWriter, req *http.Request) {
	// get vars from request to determine if environment id was specified
	vars := mux.Vars(req)
	id := vars["id"]

	// database transaction
	existingClient, err := dbGetClient(id)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}
	if existingClient == nil || existingClient.ID == "" {
		writeJSONResponse(w, http.StatusOK, successResponse{"client with specified ID does not exist"})
		return
	}
	err = dbDeleteClient(id)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}
	writeJSONResponse(w, http.StatusOK, successResponse{"OK"})
}
