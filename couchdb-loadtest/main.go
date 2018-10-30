// Simple couchdb benchmark utility (using go standard libs) used to test a cluster.
// Configurable connection settings (to test large amount of new connections) and concurrency limits.
// Useful for testing high load, high connections, and simple general performance.
// !! MADE FOR EDUCATIONAL PURPOSES; USE AT YOUR OWN RISK !!
//
// (c) 2018 George Bolo (gbolo)
// This code is licensed under MIT license
package main

import (
	"bytes"
	"crypto/sha1"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	// couchdb default endpoint settings
	default_couchdbUser     = "admin"
	default_couchdbPassword = "password"
	default_coucdbUrl       = "http://127.0.0.1:5984"
	database                = "gbolo"

	// default load test settings
	default_concurrency = 10
	default_requests    = 100

	// default http client settings
	default_maxIdleConns        = 2000
	default_maxIdleConnsPerHost = 1000
	default_disableKeepAlives   = false

	// if we should create a database
	default_createDatabase = false

	// used for test docs
	testData       = "wyKe3mXeQ9TbyO5rKtLiLn5SnMu7C6Ft0rMo1N79lSi6qY6qd2Ly8TPJ16b7lj88rUxKUJgTk2pm7yfi1b4tExs17bqhHp0pKYuXbQ1JjrckDAHKcYSfySQc4k76e8K2Ca0BHUonGWtfFEoiT0eb7pNEN1H5roT2s1LU9jRaOGNNVgLP9q5y5wL7aGw69lgfEQ2i5GGhrGmXvV7K45oUlH9NNKakTubrPSHcprFJd29DPdcnOKXE6Fk4lwEhOJDE"
	testDataUpdate = "I5XHoleO8wH3yDkQKB3VlEcZTrkrGHgAmksi9sHiRpqksGIEnFwC0s5lnq4QDuO7Amq2IfXs0exi7I9zrvDRniPakwUH7gooFPHkW5y9C2faBn5HrddtJnD6efxyLrAFbphkBQX44Oic3HO7JMi1m22ptDn6mQd6aqWCbt0dvRhj1opwv485IRq405ktBFiqruvOUfOApuWsChzHj5GS1zLZ3LHvt9vAFGMKMze82l66Wu3sk8YzpLa5UGyZnYHd"
	alphabet       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var (
	// these don't need changing
	letterRunes = []rune(alphabet)
	salt        string
	couchClient *http.Client
	totalTime   time.Time

	// load test settings
	concurrency int
	requests    int

	// couchdb endpoint
	couchdbUser, couchdbPassword, coucdbUrl string

	// http client settings
	maxIdleConns        int
	maxIdleConnsPerHost int
	disableKeepAlives   bool

	// if we should create a database
	createDatabase bool
)

// structure of couchdb doc used for testing
type Doc struct {
	// always present
	Id string `json:"_id"`
	// not needed when inserting
	Rev string `json:"_rev,omitempty"`
	// present on insert/update response
	Ok bool `json:"ok,omitempty"`

	// some data for testing
	Data string `json:"data,omitempty"`
}

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func handleError(message string, err error) {
	if err != nil {
		log.Fatalf("[%s] err: %v\n", message, err)
	}
}

func genIds(count int, salt string) (ids []string) {
	for i := 1; i <= count; i++ {
		ids = append(ids, fmt.Sprintf("%x", sha1.Sum([]byte(fmt.Sprintf("%s%s", salt, i)))))
	}
	return ids
}

func addTime(duration time.Duration) {
	totalTime = totalTime.Add(duration)
}

func createHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},

			// fabric
			ExpectContinueTimeout: 1 * time.Second,
			IdleConnTimeout:       90 * time.Second,
			MaxIdleConns:          maxIdleConns,
			MaxIdleConnsPerHost:   maxIdleConnsPerHost,
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,

			// this should disable keepalives...
			DisableKeepAlives: disableKeepAlives,
		},
	}
}

func doRequest(method string, body []byte, u *url.URL) (resHeaders http.Header, resBody []byte, resStatus int, err error) {

	// construct request
	method = strings.ToUpper(method)
	req, err := http.NewRequest(method, u.String(), bytes.NewBuffer(body))
	if err != nil {
		return
	}
	// add some headers
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	// set auth
	if len(couchdbUser) > 0 && len(couchdbPassword) > 0 {
		req.SetBasicAuth(couchdbUser, couchdbPassword)
	}

	// do the request
	start := time.Now()
	res, err := couchClient.Do(req)
	if err != nil {
		return
	}
	// calc time
	took := time.Since(start)
	fmt.Println("REQUEST took:", took)
	addTime(took)

	// close body for reuse
	defer res.Body.Close()
	resBody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	resStatus = res.StatusCode
	return

}

func getDoc(id string) (doc Doc, err error) {

	u, err := url.Parse(fmt.Sprintf("%s/%s/%s", coucdbUrl, database, id))
	if err != nil {
		return
	}

	_, body, status, err := doRequest("GET", nil, u)
	if err != nil {
		return
	}
	if status != 200 {
		err = fmt.Errorf("non-200 reponse")
		return
	}
	err = json.Unmarshal(body, &doc)
	return
}

func updateDoc(doc Doc) (docUpdated Doc, err error) {

	u, err := url.Parse(fmt.Sprintf("%s/%s/%s", coucdbUrl, database, doc.Id))
	if err != nil {
		return
	}

	reqBody, err := json.Marshal(doc)
	if err != nil {
		return
	}
	_, body, status, err := doRequest("PUT", reqBody, u)
	if err != nil {
		return
	}
	if status != 201 {
		err = fmt.Errorf("bad response code: %d", status)
		return
	}
	err = json.Unmarshal(body, &docUpdated)
	return
}

func createCouchDatabase() (err error) {
	u, err := url.Parse(fmt.Sprintf("%s/%s", coucdbUrl, database))
	if err != nil {
		return
	}

	_, body, status, err := doRequest("PUT", nil, u)
	if err != nil {
		return
	}

	fmt.Printf("CREATE DATABASE - status: %d body: %s\n", status, body)
	return
}

func runFlow(count int) {

	salt = randStringRunes(8)
	for _, id := range genIds(count, salt) {

		// PUT doc
		newDoc := Doc{Id: id, Data: testData}
		newDoc, err := updateDoc(newDoc)
		handleError("Save", err)

		// GET doc
		retrivedNewDoc, err := getDoc(id)
		handleError("Get", err)
		if retrivedNewDoc.Id != id {
			log.Fatalf("id did not equal: %s != %s\n", id, retrivedNewDoc.Id)
		}

		// UPDATE doc
		retrivedNewDoc.Data = testDataUpdate
		updatedDoc, err := updateDoc(retrivedNewDoc)
		handleError("UPDATE", err)

		// GET doc
		retrivedUpdatedDoc, err := getDoc(id)
		handleError("Get", err)

		// Validate new doc
		if retrivedUpdatedDoc.Data != testDataUpdate && retrivedUpdatedDoc.Rev != updatedDoc.Rev {
			log.Fatalf("data update not equal: %s != %s\n", testDataUpdate, retrivedUpdatedDoc.Data)
		}
	}
}

func readFlags() {
	flag.IntVar(&concurrency, "c", default_concurrency, "concurrency")
	flag.IntVar(&requests, "r", default_requests, "flows/requests")
	flag.IntVar(&maxIdleConnsPerHost, "h", default_maxIdleConnsPerHost, "maxIdleConnsPerHost")
	flag.IntVar(&maxIdleConns, "m", default_maxIdleConns, "maxIdleConns")
	flag.BoolVar(&disableKeepAlives, "k", default_disableKeepAlives, "disableKeepAlives")
	flag.BoolVar(&createDatabase, "d", default_createDatabase, "create the database")
	flag.StringVar(&coucdbUrl, "e", default_coucdbUrl, "couchdb URL")
	flag.StringVar(&couchdbUser, "u", default_couchdbUser, "username")
	flag.StringVar(&couchdbPassword, "p", default_couchdbPassword, "password")
	flag.Parse()
}

func main() {
	readFlags()
	couchClient = createHttpClient()

	// if we should create database
	if createDatabase {
		createCouchDatabase()
	}

	// do requests concurrently
	fmt.Printf("Concurrency: %d  Requests: %d\n========================================\n", concurrency, requests)
	var wg sync.WaitGroup
	wg.Add(concurrency)

	startTime := time.Now()
	totalTime = startTime
	for i := 1; i <= concurrency; i++ {
		go func() {
			defer wg.Done()
			runFlow(requests)
		}()
	}
	wg.Wait()

	// theres 4 actual http requests per flow
	fmt.Println("========================================")
	fmt.Printf("flows: %d concurrency: %d\nkeepalives: %v maxIdle: %d maxIdle/host: %d\n", requests, concurrency, !disableKeepAlives, maxIdleConns, maxIdleConnsPerHost)
	fmt.Printf("Total HTTP requests: %d\nTotal time taken by ALL HTTP requests (sum): %s\n", requests*concurrency*4, totalTime.Sub(startTime))
	fmt.Printf("Average response time per request: %vms\n", totalTime.Sub(startTime).Seconds()/float64(requests*4*concurrency)*1000)
	fmt.Printf("Total time to complete ALL flows: %s\n", time.Since(startTime))
}
