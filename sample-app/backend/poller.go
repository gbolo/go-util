package backend

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
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
	req, _ := http.NewRequest(method, reqUrl, bytes.NewBuffer(bodyBytes))
	if len(headers) > 0 {
		for header, value := range headers {
			req.Header.Set(header, value)
		}
	}
	// inject a custom user agent to easily identify this client in access logs
	req.Header.Set("User-Agent", "gbolo/sample-app/"+Version)

	// send the api request
	log.Debugf("sending http request: %s %s", method, req.URL.String())
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
	log.Debugf("response code: %d", resp.StatusCode)
	//log.Debugf("response body: %s", resp.BodyBytes)
	return
}

func getClientStatus(client *Client) (clientStatus ClientStatus) {
	clientStatus.ID = client.ID
	clientStatus.Name = client.Name
	clientStatus.URL = client.URL

	resp := doHttpCall("GET", client.URL, nil, nil)
	if resp.err != nil {
		clientStatus.Reachable = false
		clientStatus.Status = resp.err.Error()
		return
	}
	clientStatus.Reachable = true
	clientStatus.Status = fmt.Sprintf("%v", resp.StatusCode)
	return
}
