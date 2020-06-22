package backend

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"token-provider/storage"

	"github.com/gorilla/mux"
)

// @Summary Returns portal information
// @Description Returns version information
// @Tags Misc
// @Produce json
// @Success 200 {object} versionInfo
// @Router /v1/version [get]
func handlerVersion(w http.ResponseWriter, req *http.Request) {
	writeJSONResponse(w, http.StatusOK, getVersionResponse())
}

func handlerAddService(w http.ResponseWriter, req *http.Request) {
	// try to read the body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		apiResponse := errorResponse{"Bad request. Cannot read request body."}
		writeJSONResponse(w, http.StatusBadRequest, apiResponse)
		return
	}

	// try to unmarshal the body into a valid request
	var o addService
	err = json.Unmarshal(body, &o)
	if err != nil {
		apiResponse := errorResponse{"Bad request: " + err.Error()}
		writeJSONResponse(w, http.StatusBadRequest, apiResponse)
		return
	}

	// add the service
	var service *storage.Service
	if o.ID != "" {
		service = store.AddServiceWithID(o.ID, o.Description)
	} else {
		service = store.AddService(o.Description)
	}

	if service == nil {
		writeJSONResponse(w, http.StatusConflict, errorResponse{"could not add service"})
		return
	}
	writeJSONResponse(w, http.StatusAccepted, map[string]string{"id": service.ID})
}

func handlerGetServices(w http.ResponseWriter, req *http.Request) {
	services := store.ListServices()
	writeJSONResponse(w, http.StatusOK, services)
}

func handlerGetServiceByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	service := store.GetServiceByID(id)
	if service == nil {
		writeJSONResponse(w, http.StatusNotFound, errorResponse{"service ID is unknown"})
		return
	}
	writeJSONResponse(w, http.StatusOK, service)
}

func handlerGenerateApiKey(w http.ResponseWriter, req *http.Request) {
	// try to read the body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		apiResponse := errorResponse{"Bad request. Cannot read request body."}
		writeJSONResponse(w, http.StatusBadRequest, apiResponse)
		return
	}

	// try to unmarshal the body into a valid request
	var o generateAPIKey
	err = json.Unmarshal(body, &o)
	if err != nil {
		apiResponse := errorResponse{"Bad request: " + err.Error()}
		writeJSONResponse(w, http.StatusBadRequest, apiResponse)
		return
	}

	service := store.GetServiceByID(o.ServiceID)
	if service == nil {
		writeJSONResponse(w, http.StatusNotFound, errorResponse{"service ID is unknown"})
		return
	}
	rawKey, err := service.GenerateApiKey(o.Name)
	if err != nil {
		writeJSONResponse(w, http.StatusServiceUnavailable, errorResponse{"unable to generate key right now"})
		return
	}
	store.UpdateService(service)
	log.Infof("an API key was generated for service (id: %s) with prefix: %s", service.ID, rawKey.GetPrefix())
	writeJSONResponse(w, http.StatusOK, rawKey)
}

func handlerRevokeApiKey(w http.ResponseWriter, req *http.Request) {
	// try to read the body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		apiResponse := errorResponse{"Bad request. Cannot read request body."}
		writeJSONResponse(w, http.StatusBadRequest, apiResponse)
		return
	}

	// try to unmarshal the body into a valid request
	var o revokeAPIKey
	err = json.Unmarshal(body, &o)
	if err != nil {
		apiResponse := errorResponse{"Bad request: " + err.Error()}
		writeJSONResponse(w, http.StatusBadRequest, apiResponse)
		return
	}

	service := store.GetServiceByID(o.ServiceID)
	if service == nil {
		writeJSONResponse(w, http.StatusNotFound, errorResponse{"service ID is unknown"})
		return
	}
	service.RevokeApiKey(o.Prefix)
	writeJSONResponse(w, http.StatusOK, successResponse{Message: "OK"})
}

func handlerValidateKey(w http.ResponseWriter, req *http.Request) {
	// extract the key from the header
	apikeyRaw := req.Header.Get("X-API-KEY")
	if !storage.ValidateKeyFormat(apikeyRaw) {
		writeJSONResponse(w, http.StatusUnauthorized, errorResponse{"X-API-KEY is not set or invalid"})
		return
	}

	// check if api key is valid for the specified service
	vars := mux.Vars(req)
	id := vars["id"]
	if !store.ValidateApiKeyForServiceID(id, apikeyRaw) {
		writeJSONResponse(w, http.StatusUnauthorized, errorResponse{"X-API-KEY is not authorized for this service"})
		return
	}

	writeJSONResponse(w, http.StatusOK, successResponse{Message: "token is authorized"})
}

// wrapper for json responses
func writeJSONResponse(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	b, _ := json.MarshalIndent(body, "", "  ")
	w.Write(b)
}
