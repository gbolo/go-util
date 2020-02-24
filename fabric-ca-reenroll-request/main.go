package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
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

const caClientConfigContentSw = `
bccsp:
    default: SW
    sw:
        hash: SHA2
        security: 256
        filekeystore:
            keystore: msp/keystore
`

const caClientConfigContentPkcs11 = `
bccsp:
    default: PKCS11
    pkcs11:
        Library:
        Pin:
        Label:
        hash: SHA2
        security: 256
`

func main() {
	// flags to parse
	caClientBin := flag.String("b", "./fabric-ca-client", "path to fabric-ca-client binary")
	homeDir := flag.String("m", "./testdata", "base folder of msp directory")
	caProfile := flag.String("p", "", "CA profile to use")
	caName := flag.String("n", "", "CA instance name")
	debugEnabled := flag.Bool("d", false, "enable debug")
	pkcs11Enabled := flag.Bool("pkcs11", false, "enable pkcs11")
	pkcs11Library := flag.String("lib", "", "path to pkcs11 library")
	pkcs11Label := flag.String("label", "", "name of pkcs11 label/slot")
	pkcs11Pin := flag.String("pin", "", "pin for pkcs11 label/slot")
	flag.Parse()

	// confirm that the ca bin exists
	printBinVersion(*caClientBin)

	// by default use SW BCCSP config
	caClientConfigContent := caClientConfigContentSw

	// set up pkcs11 related environment variables and config if enabled
	if *pkcs11Enabled {
		caClientConfigContent = caClientConfigContentPkcs11
		fmt.Println("PKCS11 is enabled. Setting up variables")
		os.Setenv("FABRIC_CA_CLIENT_BCCSP_DEFAULT", "PKCS11")
		os.Setenv("FABRIC_CA_CLIENT_BCCSP_PKCS11_LIBRARY", *pkcs11Library)
		os.Setenv("FABRIC_CA_CLIENT_BCCSP_PKCS11_LABEL", *pkcs11Label)
		os.Setenv("FABRIC_CA_CLIENT_BCCSP_PKCS11_PIN", *pkcs11Pin)
	}

	// write fabric-ca-client config file
	configFilePath := *homeDir + "/fabric-ca-client-config.yaml"
	err := createCAClientConfig(configFilePath, caClientConfigContent)
	if err != nil {
		panic(fmt.Sprintf("unable to create config file %s: %v", configFilePath, err))
	}

	// start http server to catch response
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
		"--debug",
	)
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	// run command and ignore error (since the fabric-ca-client binary won't like the response)
	fmt.Printf("Expecting Request...\n\n")
	err = cmd.Run()
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

func createCAClientConfig(filePath, fileContent string) error {
	return ioutil.WriteFile(filePath, []byte(fileContent), 0644)
}
