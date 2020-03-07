package main

import (
	"crypto"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

const privKeyPem = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIBGI4M+crDmB/eWQj2aFxDlsiUYT1hiqW3oQA/sqQQaLoAoGCCqGSM49
AwEHoUQDQgAEm2R44sjCU5RZzkBnpCaFXakB6iBh0mqennUQBJ0g8BU7M1nxbecK
Q+hL+kF2kBxal+/fdgeOLf5W/kCkQ3O0mw==
-----END EC PRIVATE KEY-----`

func generateSignature(data []byte) (signature string, err error) {
	// decode the key, assuming it's in PEM format
	block, _ := pem.Decode([]byte(privKeyPem))
	if block == nil {
		return "", errors.New("Failed to decode PEM private key")
	}
	privKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return "", errors.New("Failed to parse ECDSA private key")
	}

	if err == nil {
		digest := sha256.Sum256(data)
		// according to https://golang.org/pkg/crypto/ecdsa/#PrivateKey.Sign
		// the SignerOpts is not actually used, its just to satisfy the crypto.Signer interface
		signatureBytes, serr := privKey.Sign(rand.Reader, digest[:], crypto.SignerOpts(crypto.SHA256))
		if serr != nil {
			panic(serr)
		}
		signature = fmt.Sprintf("%x", signatureBytes)
	}
	return
}

func main() {
	fileToSign := flag.String("input", "", "input file (file to sign)")
	outputFile := flag.String("output", "", "file to write signature to (defaults to appending .sig to input file")

	flag.Parse()
	inputFileBytes, err := ioutil.ReadFile(*fileToSign)
	if err != nil {
		fmt.Println("could not read input file:", err)
		os.Exit(1)
	}

	signature, err := generateSignature(inputFileBytes)
	if err != nil {
		fmt.Println("could not generate signature:", err)
		os.Exit(1)
	}

	if len(*outputFile) == 0 {
		*outputFile = *fileToSign + ".sig"
	}
	// TODO: consider creating a format to the signature file which describes what algo was used...
	fmt.Println("writing signature file:", *outputFile)
	err = ioutil.WriteFile(*outputFile, []byte(signature), 0644)
	if err != nil {
		fmt.Println("could not write signature file:", err)
		os.Exit(1)
	}
}
