package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func init() {
	hqClient = createHTTPClient()
}

var (
	// our global client
	hqClient *http.Client

	// agent task endpoint
	taskEndpoint = "/api/v1/task"

	// PCI compliance as of Jun 30, 2018: anything under TLS 1.1 must be disabled
	// we bump this up to TLS 1.2 so we can support best possible ciphers
	tlsMinVersion = uint16(tls.VersionTLS12)
	// allowed ciphers when in hardened mode
	// disable CBC suites (Lucky13 attack) this means TLS 1.1 can't work (no GCM)
	// only use perfect forward secrecy ciphers
	tlsCiphers = []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		// these ciphers require go 1.8+
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
	}
	// EC curve preference when in hardened mode
	// curve reference: http://safecurves.cr.yp.to/
	tlsCurvePreferences = []tls.CurveID{
		// this curve is a non-NIST curve with no NSA influence. Prefer this over all others!
		// this curve required go 1.8+
		tls.X25519,
		// These curves are provided by NIST; prefer in descending order
		tls.CurveP521,
		tls.CurveP384,
		tls.CurveP256,
	}
)

// createHTTPClient creates an http client which we can reuse
func createHTTPClient() *http.Client {
	return &http.Client{
		// http request timeout
		Timeout: 90 * time.Second,

		// transport settings
		Transport: &http.Transport{
			// TLS Config
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
				MinVersion:         tlsMinVersion,
				CipherSuites:       tlsCiphers,
				CurvePreferences:   tlsCurvePreferences,
			},

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
	jsonBytes, _ := json.Marshal(task)
	res, err := hqClient.Post(url, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return
	}

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
