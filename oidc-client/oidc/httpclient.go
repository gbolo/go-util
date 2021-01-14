package oidc

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"oidc-client/util"
	"time"
)

var httpclient *http.Client

func init() {
	// init http client
	httpclient = &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: uint16(tls.VersionTLS12),
			},
		},
	}
}

type httpResponse struct {
	StatusCode int
	BodyBytes  []byte
	err        error
}

func doHttpCall(method, reqUrl string, headers map[string]string, bodyBytes []byte) (resp httpResponse) {
	traceId := util.GenerateRandomString(6)
	req, _ := http.NewRequest(method, reqUrl, bytes.NewBuffer(bodyBytes))
	if len(headers) > 0 {
		for header, value := range headers {
			req.Header.Set(header, value)
		}
	}
	// inject a custom user agent to easily identify this client in access logs
	req.Header.Set("User-Agent", "gbolo/OIDC-TestClient")

	// send the api request
	log.Debugf("[trace-%s] sending http request: %s %s", traceId, method, req.URL.String())
	res, err := httpclient.Do(req)
	if err != nil {
		resp.err = err
		return
	}

	// read the response
	defer res.Body.Close()
	resp.StatusCode = res.StatusCode
	resp.BodyBytes, err = ioutil.ReadAll(res.Body)
	if err != nil {
		resp.err = err
		return
	}
	log.Debugf("[trace-%s] response code: %d", traceId, resp.StatusCode)
	//log.Debugf("[trace-%s] response body: %s", traceId, resp.BodyBytes)
	return
}
