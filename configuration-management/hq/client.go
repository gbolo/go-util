package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var (
	// our global client
	hqClient *http.Client

	// agent task endpoint
	taskEndpoint = "/api/v1/task"
)

// createHTTPClient creates an http client which we can reuse
func createHTTPClient() *http.Client {
	// get TLS config
	clientTLSConfig := createTLSConfig()

	return &http.Client{
		// http request timeout
		Timeout: 90 * time.Second,

		// transport settings
		Transport: &http.Transport{
			// TLS Config
			TLSClientConfig: clientTLSConfig,

			// sane timeouts
			IdleConnTimeout:       90 * time.Second,
			ResponseHeaderTimeout: 60 * time.Second,
			MaxIdleConns:          10,
			MaxIdleConnsPerHost:   5,
			DisableCompression:    true,

			// dialer timeouts
			DialContext: (&net.Dialer{
				Timeout:   90 * time.Second,
				KeepAlive: 15 * time.Second,
			}).DialContext,
		},
	}
}

// submit a task
func submitTask(url string, task task) (taskResult result, err error) {
	// create request body
	jsonBytes, _ := json.Marshal(task)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return
	}

	// set auth header if secret was specified
	if viper.GetString("secret") != "" {
		req.Header.Set("HQ-SECRET", viper.GetString("secret"))
	}

	// do request
	req.Header.Set("Content-Type", "application/json")
	res, err := hqClient.Do(req)
	if err != nil {
		return
	}

	// read response
	defer res.Body.Close()
	if res.StatusCode == 200 {
		log.Debugf("task was submitted successfully to: %v", url)
		bodyBytes, errRead := ioutil.ReadAll(res.Body)
		if errRead != nil {
			log.Errorf("error reading response body: %v", errRead)
			err = errRead
			return
		}
		err = json.Unmarshal(bodyBytes, &taskResult)
		return
	} else {
		err = fmt.Errorf("task response code was: %d", res.StatusCode)
	}
	return
}
