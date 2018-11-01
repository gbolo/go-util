# couchdb-loadtest

Simple couchdb benchmark utility (using go standard libs) used to test a cluster.
Configurable connection settings (to test large amount of new connections) and concurrency limits.
Useful for testing high load, high connections, and simple general performance.


**!! MADE FOR EDUCATIONAL PURPOSES; USE AT YOUR OWN RISK !!**

## Usage
```
# Install
go get github.com/gbolo/go-util/couchdb-loadtest

# Parameters
couchdb-loadtest --help
Usage of couchdb-loadtest:
  -c int
    	concurrency (default 10)
  -d	create the database
  -e string
    	couchdb URL (default "http://127.0.0.1:5984")
  -h int
    	maxIdleConnsPerHost (default 1000)
  -k	disableKeepAlives
  -m int
    	maxIdleConns (default 2000)
  -p string
    	password (default "password")
  -r int
    	flows/requests (default 100)
  -u string
    	username (default "admin")
```
