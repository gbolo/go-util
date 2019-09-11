package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"time"
)

var (
	listenAddress = "127.0.0.1:21700"
	requestMade   = false
)

func main() {
	// flags to parse
	caClientBin := flag.String("b", "./fabric-ca-client", "path to fabric-ca-client binary")
	homeDir := flag.String("m", "./testdata", "base folder of msp directory")
	caProfile := flag.String("p", "", "CA profile to use")
	caName := flag.String("n", "", "CA instance name")
	debugEnabled := flag.Bool("d", false, "enable debug")
	flag.Parse()

	// confirm that the ca bin exists
	printBinVersion(*caClientBin)

	// start http server
	go startHttpServer()
	time.Sleep(500 * time.Millisecond)

	// execute reenroll
	fmt.Printf("Using MSP path: %s/msp\n", *homeDir)
	cmd := exec.Command(
		*caClientBin, "reenroll",
		"--home", *homeDir,
		"--url", fmt.Sprintf("http://%s", listenAddress),
		"--enrollment.profile", *caProfile,
		"--caname", *caName,
	)
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	// run command and ignore error (since the fabric-ca-client binary won't like the response)
	fmt.Printf("Expecting Request...\n\n")
	err := cmd.Run()
	time.Sleep(200 * time.Millisecond)

	// if our server didn't receive a request then log the stdout and stderr
	if requestMade == false || *debugEnabled {
		fmt.Printf("\nDEBUG OUTPUT:\n> err: %s\n> stdOut: %s\n> stdErr: %s\n", err, out.String(), errOut.String())
	}
}

func startHttpServer() {
	// override default httpFunc
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// dump request to log
		requestDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(requestDump))

		requestMade = true
		fmt.Fprint(w, "INVALID FABRIC-CA SERVER RESPONSE")
	})

	// start the server
	fmt.Println("Starting HTTP server...")
	err := http.ListenAndServe(listenAddress, nil)
	if err != nil {
		panic(fmt.Sprintln("Error starting HTTP server:", err))
	}
}

func printBinVersion(bin string) {
	fmt.Println("using bin:", bin)
	cmd := exec.Command(bin, "version")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	fmt.Println(out.String())

	if err != nil {
		fmt.Println("Error executing fabric-ca-client binary:", err)
		os.Exit(1)
	}
}
